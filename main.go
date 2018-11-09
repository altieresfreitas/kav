package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/4ltieres/k8s-opol/pkg/admission"
	"github.com/4ltieres/k8s-opol/pkg/config"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	networkingv1 "k8s.io/api/networking/v1"
)

func getNetworkValidator(f string) *admission.NetworkAdmissionValidator {
	_, err := os.Stat(f)
	if err != nil {
		glog.Error(fmt.Errorf("Config file is missing: ", f))
	}

	yamlFile, err := os.Open(f)
	if err != nil {
		glog.Error(err)
	}

	defer yamlFile.Close()
	byteValue, _ := ioutil.ReadAll(yamlFile)
	jsonFile, err := yaml.ToJSON(byteValue)
	if err != nil {
		glog.Error(err)
	}

	policy := &admission.NetworkAdmissionValidator{}
	json.Unmarshal(jsonFile, &policy)

	return policy
}

// only allow networkpolicies with some requirements.
func admitNetworkPolicies(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	glog.V(2).Info("admitting networkpolicies")
	networkPolicyResource := metav1.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}
	if ar.Request.Resource != networkPolicyResource {
		err := fmt.Errorf("expect resource to be %s", networkPolicyResource)
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	raw := ar.Request.Object.Raw
	networkPolicy := networkingv1.NetworkPolicy{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &networkPolicy); err != nil {
		glog.Error(err)
		return toAdmissionResponse(err)
	}
	reviewResponse := v1beta1.AdmissionResponse{}

	validator := getNetworkValidator("files/config.yaml")

	ok, err := validator.IsValid(&networkPolicy)

	reviewResponse.Allowed = ok

	if !reviewResponse.Allowed {
		reviewResponse.Result = &metav1.Status{Message: strings.TrimSpace(fmt.Sprintf(err.Error()))}
	}

	return &reviewResponse
}

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
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
func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
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
		responseAdmissionReview.Response = toAdmissionResponse(err)
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

func serveNetworkPolicies(w http.ResponseWriter, r *http.Request) {
	serve(w, r, admitNetworkPolicies)
}

func main() {

	var c config.Config
	c.AddFlags()
	flag.Parse()
	http.HandleFunc("/networkpolicies", serveNetworkPolicies)

	server := &http.Server{
		Addr:      "127.0.0.1:8443",
		TLSConfig: config.ConfigTLS(c),
	}
	server.ListenAndServeTLS("", "")

}
