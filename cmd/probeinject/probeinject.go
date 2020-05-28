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
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xvzf/insight/internal/probeinject"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Logger
var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "probeinject",
		"context": "main",
	})
}

var (
	tlsCertFile = os.Getenv("TLS_CERT_FILE") // TLS Certificate path
	tlsKeyFile  = os.Getenv("TLS_KEY_FILE")  // TLS Key path
	probeImage  = os.Getenv("PROBE_IMAGE")   // Which image to use for injection
	logstash    = os.Getenv("LOGSTASH")      // Target passed to the injected network probe
)

func main() {
	// Connect to the k8s api
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}
	log.Info("Connection to Kubernetes Cluster established")

	// Probeinjector
	probeInjector := probeinject.New(clientset, "insight", corev1.Container{
		Name:  "insight-sidecar-probe",
		Image: probeImage,
		Env: []corev1.EnvVar{
			corev1.EnvVar{
				Name:  "LOGSTASH",
				Value: logstash,
			},
		},
	})
	http.HandleFunc("/inject", probeInjector.HandleWebhook)

	log.Infof("Start listening on %s", ":8443")
	err = http.ListenAndServeTLS(":8443", tlsCertFile, tlsKeyFile, nil)
	if err != nil {
		log.Fatal(err)
	}
}
