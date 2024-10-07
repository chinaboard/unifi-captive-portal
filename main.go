package main

import (
	"fmt"
	"github.com/chinaboard/unifi-captive-portal/pkg/opt"
	"github.com/chinaboard/unifi-captive-portal/pkg/unifi"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob(fmt.Sprintf("%s/*", "templates")))
)

func landingHandler(w http.ResponseWriter, r *http.Request) {
	vars := map[string]string{"Title": opt.Options.Title}
	err := templates.ExecuteTemplate(w, "landingPage", vars)
	if err != nil {
		log.Error(err.Error())
		errorHandler(w, http.StatusInternalServerError)
		return
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ap := r.URL.Query().Get("ap")
	url := r.URL.Query().Get("url")
	if id == "" || ap == "" {
		errorHandler(w, http.StatusBadRequest)
		return
	}

	client := unifi.NewClient(opt.Options.UnifiURL)
	err := client.Login(opt.Options.UnifiUsername, opt.Options.UnifiPassword)
	if err != nil {
		log.Error(err.Error())
		errorHandler(w, http.StatusInternalServerError)
		return
	}
	err = client.AuthUser(id, opt.Options.UnifiSite, opt.Options.Minutes)
	if err != nil {
		log.Error(err.Error())
		errorHandler(w, http.StatusInternalServerError)
		return
	}
	err = client.Logout()
	if err != nil {
		log.Error(err.Error())
		errorHandler(w, http.StatusInternalServerError)
		return
	}

	redirectUrl := url

	if opt.Options.RedirectUrl != "" {
		redirectUrl = opt.Options.RedirectUrl
	}

	err = templates.ExecuteTemplate(w, "thankYouPage", map[string]string{"URL": redirectUrl})
	if err != nil {
		log.Error(err)
		errorHandler(w, http.StatusInternalServerError)
		return
	}

}

func errorHandler(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	vars := map[string]string{"Title": opt.Options.Title}
	err := templates.ExecuteTemplate(w, "errorPage", vars)
	if err != nil {
		log.Error(err.Error())
	}
}

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
	log.Info("Starting up...")

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/guest/s/default/", landingHandler)
	http.HandleFunc("/auth", authHandler)
	log.Fatal(http.ListenAndServe(":80", logRequest(http.DefaultServeMux)))

}
