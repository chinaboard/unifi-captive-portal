package main

import (
	"github.com/chinaboard/unifi-captive-portal/pkg/openai"
	"github.com/chinaboard/unifi-captive-portal/pkg/options"
	"github.com/chinaboard/unifi-captive-portal/pkg/portal"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"remote_addr": r.RemoteAddr,
			"method":      r.Method,
			"url":         r.URL,
		}).Info("Request")
		handler.ServeHTTP(w, r)
	})
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("Starting up...")

	client := openai.NewClient(options.OpenAiOpt)
	portal := portal.NewPortal(options.PortalOpt)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/guest/s/default/", portal.LandingHandler)
	http.HandleFunc("/auth", portal.AuthHandler)
	http.HandleFunc("/chat", client.ChatHandler)
	log.Fatal(http.ListenAndServe(":80", logRequest(http.DefaultServeMux)))

}
