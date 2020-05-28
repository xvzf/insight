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

package kubestatestore

import (
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq" // Postgres SQL driver
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{
		"context": "kubestatestore",
	})
}

// KubeStateStore manages a stateful list of Pods mapped to their IP Addresses, Metadata and services
type KubeStateStore interface {
	HandleUpdate(e watch.Event)
}

// New creates a new pod state store with a postgres backend
func New(connString string) KubeStateStore {
	var k kubeStateStore

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	k.db = db
	log.Info("Connection to database established")

	if _, err := k.db.Exec("delete from pods"); err != nil {
		log.Error("Clear DB error:", err)
	}
	if _, err := k.db.Exec("delete from services"); err != nil {
		log.Error("Clear DB error:", err)
	}
	if _, err := k.db.Exec("delete from endpoints"); err != nil {
		log.Error("Clear DB error:", err)
	}

	return &k
}

type kubeStateStore struct {
	db *sql.DB
}

func (k *kubeStateStore) HandleUpdate(e watch.Event) {
	if e.Object == nil {
		log.Warn("event contained empty object")
		return
	}

	if pod, ok := e.Object.(*v1.Pod); ok {
		k.handlePodUpdate(e.Type, pod)
	}

	if svc, ok := e.Object.(*v1.Service); ok {
		k.handleSvcUpdate(e.Type, svc)
	}

	if endpoints, ok := e.Object.(*v1.Endpoints); ok {
		k.handleEndpointsUpdate(e.Type, endpoints)
	}
}

func (k *kubeStateStore) handlePodUpdate(event watch.EventType, pod *v1.Pod) {
	var err error
	d, _ := json.Marshal(pod)
	switch event {
	case "ADDED":
		query := "insert into pods (uid, name, namespace, ip, definition) values ($1, $2, $3, $4, $5)"
		// Not an ideal solution but it works for now. @TODO
		if pod.Status.PodIP == "" {
			_, err = k.db.Exec(query, pod.UID, pod.Name, pod.Namespace, nil, d)
		} else {
			_, err = k.db.Exec(query, pod.UID, pod.Name, pod.Namespace, pod.Status.PodIP, d)
		}
		if err != nil {
			log.WithField("kind", "pod").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": pod.Namespace,
				"name":      pod.Name,
				"ip":        pod.Status.PodIP,
			}).Info("added pod")
		}
	case "MODIFIED":
		query := "update pods set name = $2, namespace = $3, ip = $4, definition = $5 where uid = $1"
		// Not an ideal solution but it works for now. @TODO
		if pod.Status.PodIP == "" {
			_, err = k.db.Exec(query, pod.UID, pod.Name, pod.Namespace, nil, d)
		} else {
			_, err = k.db.Exec(query, pod.UID, pod.Name, pod.Namespace, pod.Status.PodIP, d)
		}
		if err != nil {
			log.WithField("kind", "pod").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": pod.Namespace,
				"name":      pod.Name,
				"ip":        pod.Status.PodIP,
			}).Info("modified pod")
		}
	case "DELETED":
		_, err = k.db.Exec("delete from pods where uid = $1", pod.UID)
		if err != nil {
			log.WithField("kind", "pod").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": pod.Namespace,
				"name":      pod.Name,
				"ip":        pod.Status.PodIP,
			}).Info("deleted pod")
		}
	default:
		log.Warn("Unknown event type" + event)
	}
}

func (k *kubeStateStore) handleSvcUpdate(event watch.EventType, svc *v1.Service) {
	var err error
	d, _ := json.Marshal(svc)
	switch event {
	case "ADDED":
		query := "insert into services (uid, name, namespace, cluster_ip, definition) values ($1, $2, $3, $4, $5)"
		// Not an ideal solution but it works for now. @TODO
		if svc.Spec.ClusterIP == "None" {
			_, err = k.db.Exec(query, svc.UID, svc.Name, svc.Namespace, nil, d)
		} else {
			_, err = k.db.Exec(query, svc.UID, svc.Name, svc.Namespace, svc.Spec.ClusterIP, d)
		}
		if err != nil {
			log.WithField("kind", "service").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": svc.Namespace,
				"name":      svc.Name,
			}).Info("added service")
		}
	case "MODIFIED":
		query := "update services set name = $2, namespace = $3, cluster_ip = $4, definition = $5 where uid = $1"
		// Not an ideal solution but it works for now. @TODO
		if svc.Spec.ClusterIP == "None" {
			_, err = k.db.Exec(query, svc.UID, svc.Name, svc.Namespace, nil, d)
		} else {
			_, err = k.db.Exec(query, svc.UID, svc.Name, svc.Namespace, svc.Spec.ClusterIP, d)
		}
		if err != nil {
			log.WithField("kind", "service").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": svc.Namespace,
				"name":      svc.Name,
			}).Info("modified service")
		}
	case "DELETED":
		_, err = k.db.Exec("delete from pods where uid = $1", svc.UID)
		if err != nil {
			log.WithField("kind", "service").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": svc.Namespace,
				"name":      svc.Name,
			}).Info("deleted service")
		}
	default:
		log.Warn("Unknown event type" + event)
	}
}

func (k *kubeStateStore) handleEndpointsUpdate(event watch.EventType, endpoints *v1.Endpoints) {
	var err error
	d, _ := json.Marshal(endpoints)
	switch event {
	case "ADDED":
		query := "insert into endpoints (uid, name, namespace, definition) values ($1, $2, $3, $4)"
		// Not an ideal solution but it works for now. @TODO
		_, err = k.db.Exec(query, endpoints.UID, endpoints.Name, endpoints.Namespace, d)
		if err != nil {
			log.WithField("kind", "endpoints").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": endpoints.Namespace,
				"name":      endpoints.Name,
			}).Info("added endpoints")
		}
	case "MODIFIED":
		query := "update endpoints set name = $2, namespace = $3, definition = $4 where uid = $1"
		_, err = k.db.Exec(query, endpoints.UID, endpoints.Name, endpoints.Namespace, d)
		if err != nil {
			log.WithField("kind", "endpoints").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": endpoints.Namespace,
				"name":      endpoints.Name,
			}).Info("modified endpoints")
		}
	case "DELETED":
		_, err = k.db.Exec("delete from endpoints where uid = $1", endpoints.UID)
		if err != nil {
			log.WithField("kind", "endpoints").Error("pq error:", err)
		} else {
			log.WithFields(logrus.Fields{
				"namespace": endpoints.Namespace,
				"name":      endpoints.Name,
			}).Info("deleted endpoints")
		}
	default:
		log.Warn("Unknown event type" + event)
	}
}
