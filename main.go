package main

import (
	"github.com/4ltieres/karepol/pkg/server"
	"github.com/golang/glog"
)

func main() {

	s := server.NewServer()

	if s.IsTLSEnable() {

		if err := s.HTTPServer.ListenAndServeTLS(s.Config.CertFile, s.Config.KeyFile); err != nil {
			glog.V(1).Infof("error %s", err)
		}

	} else {

		if err := s.HTTPServer.ListenAndServe(); err != nil {
			glog.V(1).Infof("error %s", err)
		}

	}

}
