package options

import (
	"github.com/chinaboard/unifi-captive-portal/pkg/env"
	"strconv"
)

type PortalOptions struct {
	UnifiURL      string
	UnifiUsername string
	UnifiPassword string
	UnifiSite     string
	RedirectUrl   string
	Title         string
	Minutes       string
}

type OpenAiOptions struct {
	ApiKey      string
	Domain      string
	Model       string
	Temperature float64
}

var PortalOpt *PortalOptions
var OpenAiOpt *OpenAiOptions

func init() {
	temperature, err := strconv.ParseFloat(env.Get("Temperature", "0.7"), 64)
	if err != nil {
		temperature = 0.7
	}
	OpenAiOpt = &OpenAiOptions{
		ApiKey:      env.Get("ApiKey", "sk-your-api-key"),
		Domain:      env.Get("Domain", "https://api.openai.com"),
		Model:       env.Get("Model", "gpt-3.5-turbo"),
		Temperature: temperature,
	}

	PortalOpt = &PortalOptions{
		UnifiURL:      env.Get("UnifiURL", "https://unifi:8443"),
		UnifiUsername: env.Get("UnifiUsername", "admin"),
		UnifiPassword: env.Get("UnifiPassword", "admin"),
		UnifiSite:     env.Get("UnifiSite", "default"),
		Title:         env.Get("Title", "Captive Portal"),
		RedirectUrl:   env.Get("RedirectUrl", "https://captive.apple.com/"),
		Minutes:       env.Get("Minutes", "600"),
	}

}
