package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	v1beta1 "k8s.io/api/admission/v1beta1"
	// m "github.com/alex-leonhardt/k8s-mutate-webhook/pkg/mutate"
)

func TestHandleMutateErrors(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(handleMutate))
	defer ts.Close()

	// default GET on the handle should throw an error trying to convert from empty JSON
	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.NoError(t, err)

	admReview := v1beta1.AdmissionReview{}
	assert.Errorf(t, json.Unmarshal(body, &admReview), "body: %s", string(body))

	assert.Empty(t, admReview.Response)

}
