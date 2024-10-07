package opt

import "github.com/chinaboard/unifi-captive-portal/pkg/env"

type Opt struct {
	UnifiURL      string `yaml:"unifi_url"`
	UnifiUsername string `yaml:"unifi_username"`
	UnifiPassword string `yaml:"unifi_password"`
	UnifiSite     string `yaml:"unifi_site"`
	RedirectUrl   string `yaml:"redirect_url"`
	Title         string `yaml:"title"`
	Minutes       string `yaml:"minutes"`
}

var Options = &Opt{
	UnifiURL:      env.Get("UnifiURL", "https://unifi:8443"),
	UnifiUsername: env.Get("UnifiUsername", "admin"),
	UnifiPassword: env.Get("UnifiPassword", "admin"),
	UnifiSite:     env.Get("UnifiSite", "default"),
	Title:         env.Get("Title", "Captive Portal"),
	RedirectUrl:   env.Get("RedirectUrl", "https://captive.apple.com/"),
	Minutes:       env.Get("Minutes", "600"),
}
