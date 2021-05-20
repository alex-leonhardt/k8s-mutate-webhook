package mutate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestMutatesValidRequest(t *testing.T) {
	admissionRequest := &admissionv1.AdmissionRequest{
		UID: "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
		Kind: metav1.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		},
		Resource: metav1.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "pods",
		},
		RequestKind: &metav1.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		},
		RequestResource: &metav1.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "pods",
		},
		Namespace: "yolo",
		Operation: "CREATE",
		UserInfo: authenticationv1.UserInfo{
			Username: "kubernetes-admin",
			Groups: []string{
				"system:masters",
				"system:authenticated",
			},
		},
		Object: runtime.RawExtension{
			Raw: []byte(`{
							"kind": "Pod",
							"apiVersion": "v1",
							"metadata": {
								"name": "c7m",
								"namespace": "yolo",
								"creationTimestamp": null,
								"labels": {
									"name": "c7m"
								}
							},
							"spec": {
								"volumes": [],
								"containers": [
							{
								"name": "c7m",
								"image": "primary.registry/centos:7",
								"command": [
								"/bin/bash"
							],
								"args": [
								"-c",
								"trap \"killall sleep\" TERM; trap \"kill -9 sleep\" KILL; sleep infinity"
							],
								"resources": {},
								"volumeMounts": [],
								"terminationMessagePath": "/dev/termination-log",
								"terminationMessagePolicy": "File",
								"imagePullPolicy": "IfNotPresent"
							}
							],
								"restartPolicy": "Always",
								"terminationGracePeriodSeconds": 30,
								"dnsPolicy": "ClusterFirst",
								"serviceAccountName": "default",
								"serviceAccount": "default",
								"securityContext": {},
								"schedulerName": "default-scheduler",
								"tolerations": [
							{
								"key": "node.kubernetes.io/not-ready",
								"operator": "Exists",
								"effect": "NoExecute",
								"tolerationSeconds": 300
							},
							{
								"key": "node.kubernetes.io/unreachable",
								"operator": "Exists",
								"effect": "NoExecute",
								"tolerationSeconds": 300
							}
							],
								"priority": 0,
								"enableServiceLinks": true
							},
							"status": {}
						}`),
		},
		Options: runtime.RawExtension{
			Raw: []byte(`{
							"kind": "CreateOptions",
							"apiVersion": "meta.k8s.io/v1"
						}`),
		},
	}
	response, err := Mutate(&Input{
		AdmissionRequest:  admissionRequest,
		PrimaryRegistry:   "primary.registry",
		SecondaryRegistry: "secondary.registry",
	})
	assert.NoError(t, err)

	assert.Equal(t, `[{"op":"replace","path":"/spec/containers/0/image","value":"secondary.registry/centos:7"}]`, string(response.Patch))
	assert.Contains(t, response.AuditAnnotations, "mutateme")

}
