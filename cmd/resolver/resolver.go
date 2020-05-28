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

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/sirupsen/logrus"
	"github.com/xvzf/insight/internal/resolver"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Logger
var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "resolver",
		"context": "main",
	})
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}

	log.Info(os.Getenv("INSIGHT_MEMCACHED_PORT_11211_TCP"))

	mc := memcache.New(os.Getenv("INSIGHT_MEMCACHED_PORT_11211_TCP"))

	if err := mc.Ping(); err != nil {
		log.Fatal(err)
	}

	r := resolver.New(clientset, mc)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
