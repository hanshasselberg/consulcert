package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	dst := flag.String("dst", "localhost:8300", "destination to check CA and certs from")
	flag.Parse()
	conn, err := net.Dial("tcp", *dst)
	if err != nil {
		fmt.Println("error dialing", *dst, err)
		os.Exit(1)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte{byte(3)}); err != nil {
		fmt.Println("error while writing magic TLS byte", err)
		os.Exit(1)
	}

	tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
	err = tlsConn.Handshake()
	if err != nil {
		tlsConn.Close()
		fmt.Println("error while handshake", err)
		os.Exit(1)
	}
	for _, cert := range tlsConn.ConnectionState().PeerCertificates {
		fmt.Printf("Cert Subject: %s NotBefore: %s NotAfter: %s\n", cert.Subject, cert.NotBefore, cert.NotAfter)
	}
	os.Exit(0)
}
