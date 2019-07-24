// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it, it also makes testing Mutate() kind of easy w/o need for a fake http server, etc.
package mutate

import (
	"encoding/json"
	"fmt"
	"log"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate mutates
func Mutate(body []byte) ([]byte, error) {

	log.Printf("recv: %s\n", string(body))

	responseBody := []byte{}

	var admReview *v1beta1.AdmissionReview
	var err error
	var pod *corev1.Pod

	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	ar := admReview.Request
	resp := admReview.Response

	if ar != nil {

		// Failure by default
		*resp.Result = metav1.Status{
			Status: "Failure",
		}
		// > 2 as we cater for an empty json {}
		if ar.Object.Raw != nil && len(ar.Object.Raw) > 2 {
			// get the Pod object and unmarshal it into its struct, if we cannot, we might as well stop here
			if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
				return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
			}

			// Success, of course ;)
			resp.Result = &metav1.Status{
				Status: "Success",
			}
		}

		// set response options
		resp.Allowed = true
		resp.UID = ar.UID

		pT := v1beta1.PatchTypeJSONPatch
		resp.PatchType = &pT

		// add some audit annotations, helpful to know why a object was modified, maybe (?)
		resp.AuditAnnotations = map[string]string{
			"mutateme": "yup it did it",
		}

		// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
		// tell K8S how it should modifiy it - only do this if we have already set the status to Success (the Pod json was found)
		if resp.Result.Status == "Success" {
			resp.Patch = []byte(`{ "op": "replace", "path": "/spec/containers/image", "value": "debian" }`)
		}

		// back into JSON so we can return the finished AdmissionResponse directly
		// w/o needing to convert things in the http handler
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("resp: %s\n", string(responseBody))
	return responseBody, nil
}
