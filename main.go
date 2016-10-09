package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/kargo"
)

var (
	hostname string
	httpAddr string
)

func main() {
	flag.StringVar(&httpAddr, "http", "127.0.0.1:80", "HTTP service address")
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
		link, err := kargo.Upload(kargo.UploadConfig{
			ProjectID:  "hightowerlabs",
			BucketName: "hello-universe",
			ObjectName: "hello-universe",
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(link)
		env := make(map[string]string)
		env["HELLO_UNIVERSE_TOKEN"] = os.Getenv("HELLO_UNIVERSE_TOKEN")

		dm = kargo.New("127.0.0.1:8080")

		err = dm.Create(kargo.DeploymentConfig{
			Args:      []string{"-http=0.0.0.0:80"},
			Env:       env,
			Name:      "hello-universe",
			BinaryURL: link,
		})
		if err != nil {
			log.Fatal(err)
		}

		err = dm.Logs(os.Stdout)
		if err != nil {
			log.Println(err)
		}
	} else {
		http.HandleFunc("/", httpHandler)

		go func() {
			errChan <- http.ListenAndServe(httpAddr, nil)
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
