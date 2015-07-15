package dockerevents

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
)

var (
	// DOCKER_HOST either unix://var/run/docker.sock or tcp://<ip_addr>
	endpoint = os.Getenv("DOCKER_HOST")

	// DOCKER_CERT_PATH for TLS connection to docker daemon
	certpath = os.Getenv("DOCKER_CERT_PATH")
)

func NewClient() (cli *docker.Client, err error) {
	if certpath != "" {
		ca := fmt.Sprintf("%s/ca.pem", certpath)
		cert := fmt.Sprintf("%s/cert.pem", certpath)
		key := fmt.Sprintf("%s/key.pem", certpath)
		cli, err = docker.NewTLSClient(endpoint, cert, key, ca)
	} else {
		cli, err = docker.NewClient(endpoint)
	}
	return
}
