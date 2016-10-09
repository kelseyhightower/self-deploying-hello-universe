package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/kargo"
)

var (
	cert     string
	key      string
	hostname string
	httpAddr string
)

func main() {
	flag.StringVar(&httpAddr, "http", "127.0.0.1:443", "HTTP service address")
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

	var dm *kargo.DeploymentManager
	if kargo.EnableKubernetes {
		/*
			link, err := Upload(UploadConfig{
				ProjectID:  "hightowerlabs",
				BucketName: "hello-universe",
				ObjectName: "hello-universe",
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(link)
			os.Exit(0)

		*/
		dm = kargo.New("127.0.0.1:8080")
		err := dm.Create(kargo.DeploymentConfig{
			Args: []string{"-http=0.0.0.0:443"},
			Name: "hello-universe",
		})
		if err != nil {
			log.Fatal(err)
		}
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
			if kargo.EnableKubernetes {
				err := dm.Delete()
				if err != nil {
					log.Fatal(err)
				}
			}
			os.Exit(0)
		}
	}
}
