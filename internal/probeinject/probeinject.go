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

package probeinject

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Logger
var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{
		"context": "probeinject",
	})
}

// Injector provides an HTTP-webhook interface for the kubernetes API
// and creates a json patch which will add a sidecar container
type Injector interface {
	HandleWebhook(w http.ResponseWriter, r *http.Request)
}

// RFC6902 JSON patch
type jsonPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type injector struct {
	clientSet                *kubernetes.Clientset // Kubernetes client
	watchNamespaceAnnotation string                // Annotation to look for on a NS before applying the patch
	podInjectPatch           jsonPatchOperation    // Patch operation
}

// New generates a new Injector interface providing a webhook handling kubernetes admissioncontroll calls
func New(cs *kubernetes.Clientset, annotation string, container corev1.Container) Injector {
	return &injector{
		clientSet:                cs,
		watchNamespaceAnnotation: annotation,
		podInjectPatch: jsonPatchOperation{
			Op:    "add",                // Add container
			Path:  "/spec/containers/-", // Path in the JSON object of a pod
			Value: container,            // Sidecar container to add... @TODO
		},
	}
}

// HandleWebhook handles an request from the k8s API and creates a
// json patch () injecting a sidecar proxy pushing data to elasticsearch
func (i *injector) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var body []byte
	log := log.WithField("func", "HandleWebhook")

	if r.Header.Get("Content-Type") != "application/json" {
		log.Error("invalid content type, expected application/json")
		http.Error(w, "invalid content type, expected application/json", http.StatusBadRequest)
		return
	}

	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	defer r.Body.Close()

	if len(body) == 0 {
		log.Error("empty request body")
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	// Decode AdmissionRequest
	admissionReviewRequest := admissionv1.AdmissionReview{}
	if json.Unmarshal(body, &admissionReviewRequest) != nil {
		log.Error("Failed to unmarshal AdmissionRequest object")
		http.Error(w, "failed to unmarshal AdmissionRequest", http.StatusBadRequest)
		return
	}

	// Genrate AdmissionReview including the patch for the sidecar
	resp, err := i.genAdmissionResponse(admissionReviewRequest.Request)
	if err != nil {
		log.Error(err)
		http.Error(w, "patch generation failed", http.StatusInternalServerError)
	}

	admissionReviewResponse := admissionv1.AdmissionReview{
		Response: resp,
	}

	// Debug output of what the applied patch looked like
	debug, _ := json.MarshalIndent(admissionReviewResponse, "", "  ")
	log.Info(string(debug))

	enc := json.NewEncoder(w)
	if err := enc.Encode(admissionReviewResponse); err != nil {
		log.Error(err)
		http.Error(w, "AdmissionReview encoding failed", http.StatusInternalServerError)
	} else {
		log.Info("Injected sidecar probe")
	}
}

func (i *injector) genAdmissionResponse(req *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	var pod corev1.Pod
	patchType := admissionv1.PatchTypeJSONPatch
	log := log.WithField("func", "genAdmissionResponse")

	err := json.Unmarshal(req.Object.Raw, &pod)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// More debug output :-)
	log.Info(string(req.Object.Raw))

	// Check if the namespace the pod is supposted to be deployed has probe-injection enabled
	if !i.podEnabledInjection(&pod) {
		return &admissionv1.AdmissionResponse{UID: req.UID, Allowed: true}, nil
	}

	patchBytes, err := json.Marshal([]jsonPatchOperation{i.podInjectPatch})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Build AdmissionResponse response object
	return &admissionv1.AdmissionResponse{
		UID:       req.UID,    //UID from the request needs to be passed
		Allowed:   true,       // Accept the API-Call but inject the sideacar
		Patch:     patchBytes, // Serialized array of patch operations
		PatchType: &patchType, // This needs to be a pointer, JSONpatch type is the only one supported atm
	}, nil
}

// namespaceEnabledInjection checks if the namespace has an annotation set to true
// (annotation name matching the configured state in injector)
func (i *injector) podEnabledInjection(p *corev1.Pod) bool {
	var namespace string
	log := log.WithField("func", "namespaceEnabledInjection")

	// Pod does not have a namespace assigned (yet)
	// this happens when it is created via e.g. a replicaSet
	if p.Namespace == "" {
		ors := p.GetOwnerReferences()
		// Tries to find the owner and retrieves its namespace
		for _, or := range ors {
			// If it is not a controller, skip this item
			if !*or.Controller {
				continue
			}
			// Namespace of the owner found
			if namespace != "" {
				break
			}
			switch or.Kind {
			case "ReplicaSet":
				{
					replicaSets, err := i.clientSet.AppsV1().ReplicaSets("").List(metav1.ListOptions{})
					if err != nil {
						log.Error(err)
						return false
					}
					for _, rs := range replicaSets.Items {
						if or.UID == rs.UID {
							// Found
							namespace = rs.Namespace
							break
						}
					}
				}
			case "DaemonSet":
				{
					daemonSets, err := i.clientSet.AppsV1().DaemonSets("").List(metav1.ListOptions{})
					if err != nil {
						log.Error(err)
						return false
					}
					for _, ds := range daemonSets.Items {
						if or.UID == ds.UID {
							// Found
							namespace = ds.Namespace
							break
						}
					}
				}
			case "StatefulSets":
				{
					statefulSets, err := i.clientSet.AppsV1().StatefulSets("").List(metav1.ListOptions{})
					if err != nil {
						log.Error(err)
						return false
					}
					for _, sts := range statefulSets.Items {
						if or.UID == sts.UID {
							// Found
							namespace = sts.Namespace
							break
						}
					}
				}
			}
		}
	}

	if namespace == "" {
		// @TODO add debug output
		log.Warning("Could not find namespace")
		return false
	}

	ns, err := i.clientSet.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})

	if err != nil {
		log.Error(err)
		return false
	}

	log = log.WithField("namespace", namespace)

	if err != nil {
		log.Error(err)
		return false
	}

	// Check if the annotation exists and if it is set to true; if not return false
	if v, ok := ns.GetAnnotations()[i.watchNamespaceAnnotation]; !ok || v != "true" {
		log.Info("Injection not configured/disabled")
		return false
	}

	log.Info("Injection enabled")
	return true
}
