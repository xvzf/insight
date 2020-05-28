/* 
 *  MIT License
 *  
 *  Copyright (c) 2020 Matthias Riegler <me@xvzf.tech>
 *  
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *  
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *  
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 *  
 */

package resolver

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"os/exec"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/sirupsen/logrus"
	"github.com/xvzf/insight/pkg/clusterip"
	"github.com/xvzf/insight/pkg/conntrack"
	"github.com/xvzf/insight/pkg/conntrack/parser"
	"github.com/xvzf/insight/pkg/flow/communityid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Logger
var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{
		"context": "resolver",
	})
}

// Resolver keeps track of the kubernetes cluster state and updates
// Tables in postgres
type Resolver interface {
	Run() error
	Stop()
}

// ClusterIPmap describes a mapping entry for a clusterIP -> podIP based on the communityID
type ClusterIPmap struct {
	CommunityID string `json:"community_id"` // New CommunityID
	ReplaceIP   net.IP `json:"replace_ip"`   // Actual target IP
	ReplacePort uint16 `json:"replace_port"` // Actual target port
}

type resolver struct {
	watcher   clusterip.Watcher
	clientset *kubernetes.Clientset
	memcache  *memcache.Client
	exitChan  chan struct{}
	errChan   chan error
	hasher    communityid.Hasher
}

// New creates a new Resolver
func New(clientset *kubernetes.Clientset, mc *memcache.Client) Resolver {
	return &resolver{
		watcher:   clusterip.NewWatcher(),
		clientset: clientset,
		memcache:  mc,
		exitChan:  make(chan struct{}),
		errChan:   make(chan error),
		hasher:    communityid.NewHasher(0),
	}
}

func (r *resolver) Run() error {
	// ClusterIP service watcher
	go r.clusterIPrunner()
	go r.conntrackRunner()

	err := <-r.errChan
	// Trigger exit on the other running goroutines
	close(r.exitChan)
	return err
}

func (r *resolver) Stop() {
	close(r.exitChan)
}

func (r *resolver) handleConntrackEvent(e *conntrack.Event) {
	log := log.WithField("func", "handleConntrackEvent")

	if !r.watcher.IsServiceIP(e.Entry.FlowMeta0.Dst) {
		// Conntrack entry is not a clusterIP, exit here
		return
	}

	// Calculate the replacement object
	cID := r.hasher.Hash(e.Entry.FlowMeta0)
	cm := ClusterIPmap{
		CommunityID: r.hasher.Hash(e.Entry.FlowMeta1),
		ReplaceIP:   e.Entry.FlowMeta1.Src,
		ReplacePort: e.Entry.FlowMeta1.SrcPort,
	}

	js, _ := json.Marshal(cm)

	switch e.Type {
	case conntrack.EventNew:
		{
			err := r.memcache.Set(&memcache.Item{
				Key:        cID,
				Value:      js,
				Expiration: 3600, // Default expiration time (in secondds) to be sure it gets delete
			})
			if err != nil {
				log.Error(err)
			}
			log.WithField("type", "NEW").Debug(string(js))
		}
	case conntrack.EventDestroy:
		{
			err := r.memcache.Set(&memcache.Item{
				Key:        cID,
				Value:      js,
				Expiration: 30, // Give logstash some time to process incoming flows (10s timeouts)
			})
			if err != nil {
				log.Error(err)
			}
			log.WithField("type", "DESTROY").Debug(string(js))
		}
	}
}

// conntrackRunner keeps track of conntrack connections and updates the database
func (r *resolver) conntrackRunner() {
	// @TODO move this to a netlink-native implementation
	cmd := exec.Command("conntrack", "-E")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		r.errChan <- err
		return
	}
	defer stdout.Close()

	// Create reader
	rd := bufio.NewReader(stdout)

	// Start execution
	if err := cmd.Start(); err != nil {
		r.errChan <- err
		return
	}

	for {
		// Check for exit condition
		select {
		case <-r.exitChan:
			break
		default:
		}

		// Read a line from the conntrack output
		l, err := rd.ReadString('\n')
		if err != nil {
			r.errChan <- err
			break
		}

		// Try to parse the readline
		e, err := parser.ParseConntrackEvent(l)
		// If it could be parsed, pass it on to the conntrack event handler
		if err == nil {
			go r.handleConntrackEvent(e)
		}
	}
}

// clusterIPrunner keeps track of the kubernetes service IP state
func (r *resolver) clusterIPrunner() {
	w, err := r.clientset.CoreV1().Services("").Watch(metav1.ListOptions{})

	if err != nil {
		r.errChan <- err
		return
	}
	defer w.Stop()

	for e := range w.ResultChan() {
		// Check for exit condition
		select {
		case <-r.exitChan:
			break
		default:
		}

		// Handle update
		r.watcher.HandleUpdate(e)
	}

	r.errChan <- errors.New("clusterIPrunner exited")
}
