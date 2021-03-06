package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/howler/api"
	"github.com/zalando-techmonkeys/howler/conf"
)

//Version set version information at build time
var Version = "Not set"

//Buildstamp provides a build timestamp
var Buildstamp = "Not set"

//Githash provides the current hash
var Githash = "Not set"

//serverConfig inherits Howler config
var serverConfig *conf.Config

func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
================
Example:
  %% %s
`, bin, bin)
		flag.PrintDefaults()
	}
	serverConfig = conf.New()
	serverConfig.Version = Version
	serverConfig.BuildStamp = Buildstamp
	serverConfig.GitHash = Githash
	//config from file is loaded.
	//the values will be overwritten by command line flags
	flag.BoolVar(&serverConfig.PrintVersion, "version", false, "Print version and exit")
	flag.BoolVar(&serverConfig.DebugEnabled, "debug", serverConfig.DebugEnabled, "Enable debug output")
	flag.BoolVar(&serverConfig.Oauth2Enabled, "oauth", serverConfig.Oauth2Enabled, "Enable OAuth2")
	flag.StringVar(&serverConfig.AuthURL, "oauth-authurl", serverConfig.AuthURL, "OAuth2 Auth URL")
	flag.StringVar(&serverConfig.TokenURL, "oauth-tokeninfourl", serverConfig.TokenURL, "OAuth2 Auth URL")
	flag.StringVar(&serverConfig.TLSCertfilePath, "tls-cert", serverConfig.TLSCertfilePath, "TLS Certfile")
	flag.StringVar(&serverConfig.TLSKeyfilePath, "tls-key", serverConfig.TLSKeyfilePath, "TLS Keyfile")
	flag.IntVar(&serverConfig.Port, "port", serverConfig.Port, "Listening TCP Port of the service.")
	if serverConfig.Port == 0 {
		serverConfig.Port = 1234 //default port when no option is provided
	}
	flag.DurationVar(&serverConfig.LogFlushInterval, "flush-interval", time.Second*5, "Interval to flush Logs to disk.")
}

func main() {
	flag.Parse()

	if serverConfig.PrintVersion {
		fmt.Printf("Version: %s - Build Time: %s - Git Commit Hash: %s", serverConfig.Version, serverConfig.BuildStamp, serverConfig.GitHash)
		os.Exit(0)
	}

	// default https, if cert and key are found
	var err error
	httpOnly := false
	if _, err = os.Stat(serverConfig.TLSCertfilePath); os.IsNotExist(err) {
		glog.Warningf("WARN: No Certfile found %s\n", serverConfig.TLSCertfilePath)
		httpOnly = true
	} else if _, err = os.Stat(serverConfig.TLSKeyfilePath); os.IsNotExist(err) {
		glog.Warningf("WARN: No Keyfile found %s\n", serverConfig.TLSKeyfilePath)
		httpOnly = true
	}
	var keypair tls.Certificate
	if httpOnly {
		keypair = tls.Certificate{}
	} else {
		keypair, err = tls.LoadX509KeyPair(serverConfig.TLSCertfilePath, serverConfig.TLSKeyfilePath)
		if err != nil {
			fmt.Printf("ERR: Could not load X509 KeyPair, caused by: %s\n", err)
			os.Exit(1)
		}
	}

	// configure service
	cfg := api.ServerSettings{
		Configuration: serverConfig,
		CertKeyPair:   keypair,
		Httponly:      httpOnly,
	}
	svc := api.Service{}
	svc.Run(cfg)
}
