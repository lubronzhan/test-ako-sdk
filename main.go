package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vmware/alb-sdk/go/clients"
	"github.com/vmware/alb-sdk/go/session"
)

// AviSessionTimeout is timeout for avi session
const (
	AviSessionTimeout = 60
	pageSizeMax       = "200"
	aviDefaultTenant  = "admin" // Per TKG-5862
)

func main() {
	e := os.Args[1]
	u := os.Args[2]
	p := os.Args[3]
	c := os.Args[4]
	var aviClient *clients.AviClient
	var err error

	var transport *http.Transport
	caCertPool := x509.NewCertPool()
	r, err := os.ReadFile(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read cert: %v\n", err)
		os.Exit(1)
	}

	caCertPool.AppendCertsFromPEM(r)
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}

	options := []func(*session.AviSession) error{
		session.SetPassword(p),
		session.SetTenant(aviDefaultTenant),
		session.SetControllerStatusCheckLimits(1, 1),
		session.DisableControllerStatusCheckOnFailure(true),
		session.SetTimeout(AviSessionTimeout * time.Second),
		session.SetTransport(transport),
	}

	aviClient, err = clients.NewAviClient(e, u, options...)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create API client using the credentials provided: %v\n", err)
		os.Exit(1)
	}

	apiVersion, err := aviClient.AviSession.GetControllerVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get avi controller version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("AVI Controler version: %s", apiVersion)
}
