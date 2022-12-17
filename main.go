package main

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"time"
)

func main() {
	settings, err := LoadServerFromSettings(`settings.json`)
	if err != nil {
		return
	}
	e := settings.GenerateRouter()
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(settings.CertFolder),
		HostPolicy: autocert.HostWhitelist(settings.CollectHostWhitelist()...),
	}
	serveTls := &http.Server{
		Addr:         `:443`,
		Handler:      e,
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	serve := &http.Server{
		Addr:         `:80`,
		Handler:      m.HTTPHandler(nil),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		redirectErr := serve.ListenAndServe()
		log.Printf(redirectErr.Error())
	}()
	err = serveTls.ListenAndServeTLS(``, ``)
	log.Printf(err.Error())
}
