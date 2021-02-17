package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	certFile    = "../config/client-cert.pem"
	keyFile     = "../config/client-key.pem"
	rootCAFile  = "../config/rootCA.pem"
	servers     = "nats://localhost:4222, nats://localhost:4223, nats://localhost:4224"
	streamName  = "ORDERS"
	subjectName = "ORDERS.received"
)

func main() {
	log.Print("connecting securely to cluster")
	nc, err := connect()
	noerr(err)
	defer nc.Close()

	log.Print("getting JetStream context")
	js, err := nc.JetStream()
	noerr(err)

	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Print(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subject %q", streamName, subjectName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{subjectName},
		})
		noerr(err)
	}

	log.Print("publishing an order")
	_, err = js.Publish(subjectName, []byte("one big burger"))
	noerr(err)

	log.Print("attempting to receive order")
	var order []byte
	done := make(chan bool, 1)
	js.Subscribe(subjectName, func(m *nats.Msg) {
		order = m.Data
		m.Ack()
		done <- true
	})

	select {
	case <-time.After(5 * time.Second):
		log.Fatalf("failed to get order")
	case <-done:
		log.Printf("got order: %q", order)
	}

}

// connect connects to a JetStream cluster using X509 certificates to authenticate securely.
func connect() (*nats.Conn, error) {
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing client X509 certificate/key pair: %s", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		MinVersion:   tls.VersionTLS12,
	}

	nc, err := nats.Connect(servers, nats.Secure(config), nats.RootCAs(rootCAFile))
	if err != nil {
		return nil, fmt.Errorf("Got an error on Connect with Secure Options: %s\n", err)
	}
	return nc, nil
}

func noerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
