package main

import (
	"crypto/tls"
	"fmt"
	"forum/database"
	"forum/handlers"
	"net/http"
)

func main() {
	const (
		CertFilePath = "./tls/cert.pem"
		KeyFilePath  = "./tls/key.pem"
	)

	database.InitDB()

	rateLimiter := handlers.NewRateLimiter()

	h := routes()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":8080",
		Handler:   rateLimiter.LimitMiddleware(h),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running at https://localhost:8080")
	err := server.ListenAndServeTLS(CertFilePath, KeyFilePath)
	if err != nil {
		fmt.Println("Error starting HTTPS server:", err)
	}
}
