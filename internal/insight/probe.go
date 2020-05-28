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

package insight

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/sirupsen/logrus"
	"github.com/xvzf/insight/pkg/capture"
	"github.com/xvzf/insight/pkg/flow/container"
	"github.com/xvzf/insight/pkg/insight"
)

// Logger
var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{
		"context": "insight",
		"int":     "probe",
	})
}

// Probe is a task runner destinated to manage flow containers and fill them
type Probe interface {
	Run() error
	Stop()
}

type probe struct {
	sampleTime     time.Duration         // How often should we create a new flow container
	logstash       string                // Logstash target
	capture        capture.Capturer      // Capture object
	container      container.Container   // Flow container
	containerMutex sync.Mutex            // goroutine safe :-)
	dumpChan       chan []*insight.Event // Dump channel
	exitChan       chan struct{}         // Exit channel
	errChan        chan error            // Error channel
}

// NewProbe creates a new probe object
func NewProbe(c capture.Capturer, st time.Duration, l string) Probe {
	p := &probe{
		capture:    c,
		sampleTime: st,
		logstash:   l,
		container:  container.New(),
		dumpChan:   make(chan []*insight.Event, 10),
		exitChan:   make(chan struct{}),
		errChan:    make(chan error),
	}

	return p
}

func (p *probe) newContainer() {
	p.containerMutex.Lock()
	defer p.containerMutex.Unlock()
	log.Info("Creating new flow container")
	// convert to events & transmit
	select {
	case p.dumpChan <- insight.NewFromFlows(p.container.Dump()):
	default:
		p.errChan <- errors.New("Buffer full, dropping flows")
	}

	p.container = container.New()
}

func (p *probe) addEvents(events []*insight.Event) {
}

func (p *probe) handlePacket(gp gopacket.Packet) {
	s, err := capture.NewSample(gp)
	if err != nil {
		log.Warn(err)
		return
	}
	p.container.Add(s)
}

func (p *probe) captureRunner() {
	for gp := range p.capture.Packets() {
		select {
		case <-p.errChan:
			return
		case <-p.exitChan:
			return
		default:
		}
		go p.handlePacket(gp)
	}
	p.errChan <- errors.New("captureRunner exited")
}

func (p *probe) dumpRunner() {
	for {
		select {
		case <-p.errChan:
			return
		case <-p.exitChan:
			return
		case events := <-p.dumpChan:
			js, err := json.Marshal(events)
			if err != nil {
				p.errChan <- err
			}
			resp, err := http.Post(p.logstash, "application/json", bytes.NewBuffer(js))
			if err != nil || resp.StatusCode != 200 {
				log.Errorf("Failed to dump buffer, %d flow records lost", len(events))
			}
			log.WithField("container_size", len(events)).Info("Dumped container")
		}
	}
}

func (p *probe) containerRenewRunner() {
	for {
		select {
		case <-p.errChan:
			return
		case <-p.exitChan:
			return
		default:
			time.Sleep(p.sampleTime)
		}
		p.newContainer()
	}
}

func (p *probe) Run() error {
	go p.captureRunner()
	go p.containerRenewRunner()
	go p.dumpRunner()
	return <-p.errChan
}

func (p *probe) Stop() {
	close(p.exitChan)
	p.errChan <- nil
}
