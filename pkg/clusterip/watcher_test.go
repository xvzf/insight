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

import "testing"

import (
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"net"
)

var (
	eventstream []watch.Event = []watch.Event{
		watch.Event{
			Type: "ADDED",
			Object: &corev1.Service{
				Spec: corev1.ServiceSpec{ClusterIP: "10.0.0.1"},
			},
		},
		watch.Event{
			Type: "ADDED",
			Object: &corev1.Service{
				Spec: corev1.ServiceSpec{ClusterIP: "10.0.0.2"},
			},
		},
		watch.Event{
			Type: "ADDED",
			Object: &corev1.Service{
				Spec: corev1.ServiceSpec{ClusterIP: "10.0.0.3"},
			},
		},
		watch.Event{
			Type: "DELETED",
			Object: &corev1.Service{
				Spec: corev1.ServiceSpec{ClusterIP: "10.0.0.3"},
			},
		},
		watch.Event{
			Type:   "ADDED",
			Object: &corev1.Pod{},
		},
		watch.Event{
			Type: "ADDED",
			Object: &corev1.Service{
				Spec: corev1.ServiceSpec{ClusterIP: "10.0.0.4"},
			},
		},
		watch.Event{
			Type: "ADDED",
			Object: &corev1.Service{
				Spec: corev1.ServiceSpec{ClusterIP: "None"},
			},
		},
	}

	mapAfterEvents map[string]bool = map[string]bool{
		"10.0.0.1": true,
		"10.0.0.2": true,
		"10.0.0.4": true,
	}

	clusterIPs []net.IP = []net.IP{
		net.ParseIP("10.0.0.1"),
		net.ParseIP("10.0.0.2"),
		net.ParseIP("10.0.0.4"),
	}

	nonClusterIPs []net.IP = []net.IP{
		net.ParseIP("10.0.0.3"),
		net.ParseIP("10.1.0.1"),
		net.ParseIP("10.1.0.2"),
		net.ParseIP("10.1.0.4"),
	}
)

func TestNewWatcher(t *testing.T) {
	w := NewWatcher()
	rawW, ok := w.(*watcher)
	if !ok {
		t.Error("New returns invalid object type, expected *watcher")
		return
	}

	if rawW.data == nil {
		t.Error("Datastructure not initialized")
	}
}

func buildWatcher() Watcher {
	w := NewWatcher()
	for _, e := range eventstream {
		w.HandleUpdate(e)
	}
	return w
}

func TestHandleUpdate(t *testing.T) {
	w := buildWatcher()

	rawW, _ := w.(*watcher)

	if !cmp.Equal(rawW.data, mapAfterEvents) {
		t.Error(cmp.Diff(rawW.data, mapAfterEvents))
	}
}

func TestIsServiceIP(t *testing.T) {
	w := buildWatcher()

	for _, ip := range clusterIPs {
		if !w.IsServiceIP(ip) {
			t.Errorf("%s should be a cluster IP", ip)
		}
	}

	for _, ip := range nonClusterIPs {
		if w.IsServiceIP(ip) {
			t.Errorf("%s should not be a cluster IP", ip)
		}
	}
}
