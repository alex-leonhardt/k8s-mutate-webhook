package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/loadsmart/k8s-mutate-webhook/pkg/mutate"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
)

func mutateHandler(origRegistry, newRegistry string, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read the body / request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			sendError(err, w, logger)
			return
		}

		// unmarshal request into AdmissionReview struct
		admReview := admissionv1.AdmissionReview{}
		if err := json.Unmarshal(body, &admReview); err != nil {
			sendError(fmt.Errorf("unmarshaling request failed with %s", err), w, logger)
			return
		}

		// mutate the request
		admReview.Response, err = mutate.Mutate(&mutate.Input{
			AdmissionRequest:  admReview.Request,
			PrimaryRegistry:   origRegistry,
			SecondaryRegistry: newRegistry,
		})
		if err != nil {
			sendError(err, w, logger)
			return
		}

		responseBody, err := json.Marshal(admReview)
		if err != nil {
			sendError(err, w, logger)
			return
		}

		// and write it back
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseBody)
	})
}

func sendError(err error, w http.ResponseWriter, logger *zap.Logger) {
	logger.Error(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprintf(w, "%s", err)
}

func NewMutatePodsHandler(r *mux.Router, origRegistry, newRegistry string, logger *zap.Logger) {
	r.Handle("/mutate/pods", mutateHandler(origRegistry, newRegistry, logger)).Methods("POST").Headers("Content-Type", "application/json")
}
