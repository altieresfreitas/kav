package server

import (
	"flag"
	"net/http"

	"github.com/4ltieres/k8s-opol/pkg/admission"

	"github.com/4ltieres/k8s-opol/pkg/config"
)

// Server is an abstraction
type Server struct {
	HTTPServer *http.Server
	Config     config.Config
	Validator  *admission.NetworkAdmissionValidator
}

func (s *Server) serveNetworkPolicies(w http.ResponseWriter, r *http.Request) {
	s.serve(w, r, s.admitNetworkPolicies)
}

//IsTLSEnable Check if tls configs was passed
func (s *Server) IsTLSEnable() bool {
	if s.Config.KeyFile != "" && s.Config.CertFile != "" {

		return true

	}

	return false
}

// NewServer return a server
func NewServer() *Server {

	c := config.Config{}
	c.AddFlags()
	flag.Parse()

	s := &Server{
		HTTPServer: &http.Server{
			Addr: c.ListenAddress,
		},
		Config: c}

	http.HandleFunc("/networkpolicies", s.serveNetworkPolicies)
	return s
}
