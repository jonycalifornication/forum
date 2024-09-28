package main

import (
	"crypto/tls"
	"fmt"
	"forum/database"
	"forum/internal"
	"log"
	"net/http"
)

func main() {
	const (
		CertFilePath = "./tls/cert.pem"
		KeyFilePath  = "./tls/key.pem"
	)

	database.InitDB()

	_, err := internal.LoadConfig("config.json")
	if err != nil {
		log.Println(err)
	}

	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		log.Println(err)
	}

	h := routes()

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		MinVersion:   tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":8080",
		Handler:   h,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running at http://localhost:8080")
	err2 := server.ListenAndServeTLS("", "")
	if err2 != nil {
		fmt.Println("Error starting HTTPS server:", err2)
	}
}
