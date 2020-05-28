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

package main

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/xvzf/insight/internal/kubestatestore"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{
		"context": "kubeagent",
	})
}

func main() {
	log.Info("Starting up KubeAgent")

	// @TODO connect out of container
	// var kubeconfig = "/Users/xvzf/.kube/config"
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}
	log.Info("Connection to Kubernetetes Cluster established")

	podWatcher, err := clientset.CoreV1().Pods("").Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to watch pod API", err)
	}

	svcWatcher, err := clientset.CoreV1().Services("").Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to watch service API", err)
	}

	endpointsWatcher, err := clientset.CoreV1().Endpoints("").Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to watch endpoint API", err)
	}

	var wg sync.WaitGroup
	store := kubestatestore.New(os.Getenv("CONN_STRING"))
	// Create a goroutine for the pod & service watcher
	for _, watcher := range []watch.Interface{podWatcher, svcWatcher, endpointsWatcher} {
		wg.Add(1)
		// Capture incoming events and pass them to the PodStateStore
		go func(w watch.Interface) {
			defer wg.Done()
			for e := range w.ResultChan() {
				// Pass event to the store
				go store.HandleUpdate(e)
			}
		}(watcher)
	}
	wg.Wait()
}
