package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/hello-universe/kubernetes"
)

var (
	cert             string
	key              string
	expose           bool
	replicas         int
	hostname         string
	httpAddr         string
	enableKubernetes bool
)

func main() {
	flag.StringVar(&httpAddr, "http", "127.0.0.1:443", "HTTP service address")
	flag.BoolVar(&expose, "expose", false, "Create a Kubernetes")
	flag.BoolVar(&enableKubernetes, "kubernetes", false, "Deploy to Kubernetes.")
	flag.IntVar(&replicas, "replicas", 1, "Number of replicas")
	flag.StringVar(&cert, "cert", "/etc/hello-universe/tls.crt", "TLS certificate path")
	flag.StringVar(&key, "key", "/etc/hello-universe/tls.key", "TLS private key path")
	flag.Parse()

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initializing hello-universe ...")
	errChan := make(chan error, 10)

	if enableKubernetes {
		if err := kubernetes.CreateSecret(key, cert); err != nil {
			log.Fatal(err)
		}
		if err := kubernetes.CreateDeployment(); err != nil {
			log.Fatal(err)
		}
		go kubernetes.Logs()
	} else {
		http.HandleFunc("/", httpHandler)

		go func() {
			errChan <- http.ListenAndServeTLS(httpAddr, cert, key, nil)
		}()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			if enableKubernetes {
				kubernetes.DeleteDeployment()
				kubernetes.DeleteSecret()
			}
			os.Exit(0)
		}
	}
}
