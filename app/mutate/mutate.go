// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it, it also makes testing Mutate() kind of easy w/o need for a fake http server, etc.
package mutate

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// jsonPatch helps marshal the patch operation as JSON
type jsonPatch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// Mutate mutates
func Mutate(body []byte, labels map[string]string, annotations map[string]string) ([]byte, error) {
	log.Printf("recv: %s\n", string(body)) // untested section

	// unmarshal request into AdmissionReview struct
	admReview := v1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	resp := v1.AdmissionResponse{}

	if admReview.Request == nil {
		return []byte{}, nil
	}

	// get the Pod object and unmarshal it into its struct, if we cannot, we might as well stop here
	var pod *corev1.Pod
	if err := json.Unmarshal(admReview.Request.Object.Raw, &pod); err != nil {
		return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
	}
	// set response options
	resp.Allowed = true
	resp.UID = admReview.Request.UID
	pT := v1.PatchTypeJSONPatch
	resp.PatchType = &pT // it's annoying that this needs to be a pointer as you cannot give a pointer to a constant?

	// add some audit annotations, helpful to know why a object was modified, maybe (?)
	resp.AuditAnnotations = map[string]string{
		"k8s-pod-labeler": fmt.Sprintf("adding (%d) labels and (%d) annotations", len(labels), len(annotations)),
	}

	// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
	// tell K8S how it should modifiy it
	patches := make([]jsonPatch, 0, len(annotations)+len(labels))

	// add each new label
	hasLabels := len(pod.ObjectMeta.Labels) > 0
	for label, value := range labels {
		log.Printf("adding label %s", label)
		if _, ok := pod.ObjectMeta.Labels[label]; !ok {
			var patch jsonPatch
			// if first label, add the initial label object
			if !hasLabels {
				patch = jsonPatch{
					Op:    "add",
					Path:  "/metadata/labels",
					Value: map[string]string{label: value},
				}
			} else {
				patch = jsonPatch{
					Op:    "add",
					Path:  fmt.Sprintf("/metadata/labels/%s", jsonPointersEncode(label)),
					Value: value,
				}
			}
			patches = append(patches, patch)
			hasLabels = true
		} else {
			log.Printf("skipping label, already exists, %s", label)
		}
	}

	// add each new annotation
	hasAnnotations := len(pod.ObjectMeta.Annotations) > 0
	for annotation, value := range annotations {
		log.Printf("adding annotation %s", annotation)
		if _, ok := pod.ObjectMeta.Annotations[annotation]; !ok {
			var patch jsonPatch
			if !hasAnnotations {
				patch = jsonPatch{
					Op:    "add",
					Path:  "/metadata/annotations",
					Value: map[string]string{annotation: value},
				}
			} else {
				patch = jsonPatch{
					Op:    "add",
					Path:  fmt.Sprintf("/metadata/annotations/%s", jsonPointersEncode(annotation)),
					Value: value,
				}
			}
			patches = append(patches, patch)
			hasAnnotations = true
		} else {
			log.Printf("skipping annotation, already exists, %s", annotation)
		}
	}

	// parse the []map into JSON
	respPatch, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}
	resp.Patch = respPatch

	// Success, of course ;)
	resp.Result = &metav1.Status{
		Status: metav1.StatusSuccess,
	}

	admReview.Response = &resp
	// back into JSON so we can return the finished AdmissionReview w/ Response directly
	// w/o needing to convert things in the http handler
	responseBody, err := json.Marshal(admReview)
	if err != nil {
		return nil, err // untested section
	}

	return responseBody, nil
}

// jsonPointersEncode implements jsonpath encoding
// https://datatracker.ietf.org/doc/html/rfc6901#section-3
func jsonPointersEncode(in string) string {
	out := in
	out = strings.ReplaceAll(out, "~", "~0")
	out = strings.ReplaceAll(out, "/", "~1")
	return out
}
