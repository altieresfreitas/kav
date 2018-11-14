package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/4ltieres/k8s-opol/pkg/admission"
	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	networkingv1 "k8s.io/api/networking/v1"
)

// only allow networkpolicies with some requirements.
func (s *Server) admitNetworkPolicies(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {

	glog.V(2).Info("admitting networkpolicies")
	var group, version, resource string

	switch ar.Request.Resource {
	case metav1.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}:

		group = "networking.k8s.io"
		version = "v1"
		resource = "networkpolicies"

	case metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "networkpolicies"}:
		group = "extensions"
		version = "v1beta1"
		resource = "networkpolicies"

	}

	networkPolicyResource := metav1.GroupVersionResource{Group: group, Version: version, Resource: resource}

	if ar.Request.Resource != networkPolicyResource {
		err := fmt.Errorf("expect resource to be %s ", networkPolicyResource)
		glog.Error(err)
		return s.toAdmissionResponse(err)
	}

	raw := ar.Request.Object.Raw
	networkPolicy := networkingv1.NetworkPolicy{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &networkPolicy); err != nil {
		glog.Error(err)
		return s.toAdmissionResponse(err)
	}
	reviewResponse := v1beta1.AdmissionResponse{}

	validator := admission.NewAdmissionValidator(s.Config.ConfigFile)

	ok, err := validator.IsValid(&networkPolicy)

	reviewResponse.Allowed = ok

	if !reviewResponse.Allowed {
		reviewResponse.Result = &metav1.Status{Message: strings.TrimSpace(fmt.Sprintf(err.Error()))}
	}

	return &reviewResponse
}

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func (s *Server) toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// serve handles the http portion of a request prior to handing to an admit
// function
func (s *Server) serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	glog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1beta1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1beta1.AdmissionReview{}

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		glog.Error(err)
		responseAdmissionReview.Response = s.toAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = admit(requestedAdmissionReview)
	}

	// Return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	glog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		glog.Error(err)
	}
	if _, err := w.Write(respBytes); err != nil {
		glog.Error(err)
	}
}
