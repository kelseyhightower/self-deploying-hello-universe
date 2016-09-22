package kubernetes

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getPods() ([]string, error) {
	time.Sleep(10 * time.Second)
	cmd := exec.Command("kubectl", "get", "pods", "-l", "app=hello-universe", "-o", "name")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
	}

	pods := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		pods = append(pods, scanner.Text())
	}
	return pods, nil
}

func Logs() error {
	pods, err := getPods()
	if err != nil {
		return err
	}
	for _, pod := range pods {
		pod := pod
		go func() {
			cmd := exec.Command("kubectl", "logs", "-f", pod)
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Println(err)
			}
			if err := cmd.Start(); err != nil {
				log.Println(err)
			}
			go io.Copy(os.Stdout, stdout)
		}()
	}
	return nil
}

func CreateSecret(key, cert string) error {
	log.Println("Creating hello-universe secret...")
	cmd := exec.Command("kubectl", "create", "secret", "tls", "hello-universe",
		"--key", key, "--cert", cert)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
	}
	return err
}

func DeleteSecret() error {
	log.Println("Deleting hello-universe secret...")
	cmd := exec.Command("kubectl", "delete", "secret", "hello-universe")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
	}
	return err
}

func CreateDeployment() error {
	log.Println("Creating hello-universe deployment...")
	cmd := exec.Command("kubectl", "create", "-f", "-")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = strings.NewReader(deploymentConfig)

	err := cmd.Run()
	if err != nil {
		log.Println(out.String)
	}
	return err
}

func DeleteDeployment() error {
	log.Println("Deleting hello-universe deployment...")
	cmd := exec.Command("kubectl", "delete", "deployment", "hello-universe")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
	}
	return err
}
