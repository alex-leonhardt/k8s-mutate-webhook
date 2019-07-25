package mutate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestMutateJSON(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "Pod",
				"apiVersion": "v1",
				"metadata": {
					"name": "c7m",
					"namespace": "yolo",
					"creationTimestamp": null,
					"labels": {
						"name": "c7m"
					},
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"c7m\"},\"name\":\"c7m\",\"namespace\":\"yolo\"},\"spec\":{\"containers\":[{\"args\":[\"-c\",\"trap \\\"killall sleep\\\" TERM; trap \\\"kill -9 sleep\\\" KILL; sleep infinity\"],\"command\":[\"/bin/bash\"],\"image\":\"centos:7\",\"name\":\"c7m\"}]}}\n"
					}
				},
				"spec": {
					"volumes": [
						{
							"name": "default-token-5z7xl",
							"secret": {
								"secretName": "default-token-5z7xl"
							}
						}
					],
					"containers": [
						{
							"name": "c7m",
							"image": "centos:7",
							"command": [
								"/bin/bash"
							],
							"args": [
								"-c",
								"trap \"killall sleep\" TERM; trap \"kill -9 sleep\" KILL; sleep infinity"
							],
							"resources": {},
							"volumeMounts": [
								{
									"name": "default-token-5z7xl",
									"readOnly": true,
									"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
								}
							],
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
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	response, err := Mutate([]byte(rawJSON))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"replace","path":"/spec/containers/0/image","value":"debian"}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")

}
