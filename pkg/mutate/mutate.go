package mutate

import (
	"encoding/json"
	"fmt"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Input struct {
	AdmissionRequest *admissionv1.AdmissionRequest
	PrimaryRegistry string
	SecondaryRegistry string
	Verbose bool
}

func Mutate(input *Input) (*admissionv1.AdmissionResponse, error) {
	var pod *corev1.Pod

	admissionResponse := &admissionv1.AdmissionResponse{}

	if input.AdmissionRequest != nil {

		// get the Pod object and unmarshal it into its struct, if we cannot, we might as well stop here
		if err := json.Unmarshal(input.AdmissionRequest.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("unable to unmarshal pod json object %v", err)
		}
		// set response options
		admissionResponse.Allowed = true
		admissionResponse.UID = input.AdmissionRequest.UID
		pT := admissionv1.PatchTypeJSONPatch
		admissionResponse.PatchType = &pT // it's annoying that this needs to be a pointer as you cannot give a pointer to a constant?

		// add some audit annotations, helpful to know why a object was modified, maybe (?)
		admissionResponse.AuditAnnotations = map[string]string{
			"mutateme": "yup it did it",
		}

		// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
		// tell K8S how it should modifiy it
		var p []map[string]string
		for i, container := range pod.Spec.Containers {
			if strings.HasPrefix(container.Image, input.PrimaryRegistry) {
				newImage := input.SecondaryRegistry + strings.TrimPrefix(container.Image, input.PrimaryRegistry)
				patch := map[string]string{
					"op":    "replace",
					"path":  fmt.Sprintf("/spec/containers/%d/image", i),
					"value": newImage,
				}
				p = append(p, patch)
			}
		}
		// parse the []map into JSON
		admissionResponse.Patch, _ = json.Marshal(p)

		// Success, of course ;)
		admissionResponse.Result = &metav1.Status{
			Status: "Success",
		}

	}

	return admissionResponse, nil
}
