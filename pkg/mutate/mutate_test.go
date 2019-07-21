package mutate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestMutate(t *testing.T) {

	p := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Labels: map[string]string{
				"Name": "FakePod",
			},
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				corev1.Container{
					Name:  "c1",
					Image: "yolo",
				},
			},
		},
	}
	pj, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("failed to marshal %v to json with error %s", p, err)
	}

	AdmissionRequestPod := v1beta1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind: "AdmissionReview",
		},
		Request: &v1beta1.AdmissionRequest{
			UID: "e911857d-c318-11e8-bbad-025000000001",
			Kind: metav1.GroupVersionKind{
				Kind:    "Pod",
				Version: "v1",
			},
			Operation: "CREATE",
			Object: runtime.RawExtension{
				Raw: pj,
			},
		},
	}

	request, err := json.Marshal(AdmissionRequestPod.Request)
	if err != nil {
		t.Fatalf("failed to create AdmissionRequest json with error %s", err)
	}

	response, err := Mutate([]byte(request))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(request), err)
	}

	r := v1beta1.AdmissionResponse{}
	if err := json.Unmarshal(response, &r); err != nil {
		t.Errorf("failed to unmarshal %s with error %s", response, err)
	}

	assert.Equal(t, string(r.Patch), `{ "op": "replace", "path": "/spec/containers/image", "value": "debian" }`)
	assert.Equal(t, AdmissionRequestPod.Request.UID, r.UID)
	assert.Contains(t, r.AuditAnnotations, "mutateme")

}

func TestMutateNoRawObject(t *testing.T) {

	AdmissionRequestPod := v1beta1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind: "AdmissionReview",
		},
		Request: &v1beta1.AdmissionRequest{
			UID: "e911857d-c318-11e8-bbad-025000000002",
			Kind: metav1.GroupVersionKind{
				Kind:    "Pod",
				Version: "v1",
			},
			Operation: "CREATE",
			Object: runtime.RawExtension{
				Raw: []byte(`{}`),
			},
		},
	}

	request, err := json.Marshal(AdmissionRequestPod.Request)
	if err != nil {
		t.Fatalf("failed to create AdmissionRequest json with error %s", err)
	}

	response, err := Mutate([]byte(request))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(request), err)
	}

	r := v1beta1.AdmissionResponse{}
	if err := json.Unmarshal(response, &r); err != nil {
		t.Errorf("failed to unmarshal %s with error %s", response, err)
	}

	assert.Equal(t, &metav1.Status{Status: "Failure"}, r.Result)
	assert.NotEqual(t, `{ "op": "replace", "path": "/spec/containers/image", "value": "debian" }`, string(r.Patch))
	assert.Equal(t, AdmissionRequestPod.Request.UID, r.UID)
	assert.Contains(t, r.AuditAnnotations, "mutateme")

}
