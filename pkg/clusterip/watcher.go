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

package clusterip

import (
	"net"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watcher contains a kubernetes service-ip state and provides an API
// to pass k8s events and check whether an IP is a k8s ClusterIP
type Watcher interface {
	HandleUpdate(e watch.Event)
	IsServiceIP(ip net.IP) bool
}

type watcher struct {
	sync.Mutex
	data map[string]bool
}

// NewWatcher creates a new Watcher
func NewWatcher() Watcher {
	return &watcher{
		data: make(map[string]bool),
	}
}

// HandleUpdate handles API Server requests
func (w *watcher) HandleUpdate(e watch.Event) {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()
	svc, ok := e.Object.(*corev1.Service)

	// Not a service
	if !ok {
		return
	}

	// Not of type clusterIP
	if svc.Spec.ClusterIP == "None" {
		return
	}

	switch e.Type {
	case "ADDED":
		w.data[svc.Spec.ClusterIP] = true
	case "DELETED":
		delete(w.data, svc.Spec.ClusterIP)
	}
}

// IsServiceIP checks if the processed IP belongs to a service
func (w *watcher) IsServiceIP(ip net.IP) bool {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	v, ok := w.data[ip.String()]
	return (v && ok)
}
