package config

import (
	"crypto/tls"
	"flag"

	"github.com/golang/glog"
)

// Config contains the server (the webhook) cert and key.
type Config struct {
	CertFile      string
	KeyFile       string
	ConfigFile    string
	ListenAddress string
}

// AddFlags parse flags
func (c *Config) AddFlags() {

	flag.StringVar(&c.CertFile, "tls-cert-file", c.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&c.KeyFile, "tls-private-key-file", c.KeyFile, ""+
		"File containing the default x509 private key matching --tls-cert-file.")
	flag.StringVar(&c.ConfigFile, "config-file", c.ConfigFile, ""+
		"File containing validation rules --config-file.")
	flag.StringVar(&c.ListenAddress, "listen-address", "0.0.0.0:443", ""+
		"File containing validation rules --listen-address.")

}

// ConfigTLS return a tls object
func ConfigTLS(c Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		glog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}
