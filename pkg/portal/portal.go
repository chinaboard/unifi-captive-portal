package portal

import (
	"fmt"
	"github.com/chinaboard/unifi-captive-portal/pkg/options"
	"github.com/chinaboard/unifi-captive-portal/pkg/unifi"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type Portal struct {
	opt       *options.PortalOptions
	templates *template.Template
}

func NewPortal(opt *options.PortalOptions) *Portal {
	return &Portal{
		opt:       opt,
		templates: template.Must(template.ParseGlob(fmt.Sprintf("%s/*", "templates"))),
	}
}

func (p *Portal) LandingHandler(w http.ResponseWriter, r *http.Request) {
	vars := map[string]string{"Title": p.opt.Title}
	err := p.templates.ExecuteTemplate(w, "landingPage", vars)
	if err != nil {
		log.Error(err.Error())
		p.errorHandler(w, http.StatusInternalServerError)
		return
	}
}

func (p *Portal) AuthHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ap := r.URL.Query().Get("ap")
	url := r.URL.Query().Get("url")
	if id == "" || ap == "" {
		p.errorHandler(w, http.StatusBadRequest)
		return
	}

	client, err := unifi.NewClient(p.opt.UnifiURL)
	if err != nil {
		log.Error(err.Error())
		p.errorHandler(w, http.StatusInternalServerError)
		return
	}
	err = client.Login(p.opt.UnifiUsername, p.opt.UnifiPassword)
	if err != nil {
		log.Error(err.Error())
		p.errorHandler(w, http.StatusInternalServerError)
		return
	}
	err = client.AuthUser(id, p.opt.UnifiSite, p.opt.Minutes)
	if err != nil {
		log.Error(err.Error())
		p.errorHandler(w, http.StatusInternalServerError)
		return
	}
	err = client.Logout()
	if err != nil {
		log.Error(err.Error())
		p.errorHandler(w, http.StatusInternalServerError)
		return
	}

	redirectUrl := url

	if p.opt.RedirectUrl != "" {
		redirectUrl = p.opt.RedirectUrl
	}

	err = p.templates.ExecuteTemplate(w, "thankYouPage", map[string]string{"URL": redirectUrl})
	if err != nil {
		log.Error(err)
		p.errorHandler(w, http.StatusInternalServerError)
		return
	}

}

func (p *Portal) errorHandler(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	vars := map[string]string{"Title": p.opt.Title}
	err := p.templates.ExecuteTemplate(w, "errorPage", vars)
	if err != nil {
		log.Error(err.Error())
	}
}
