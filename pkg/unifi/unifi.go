package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Client struct {
	client  *http.Client
	baseUrl string
}

func NewClient(url string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		client: &http.Client{Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}},
		baseUrl: strings.TrimSuffix(url, "/"),
	}
}

func (u *Client) do(path string, body any) (int, error) {
	json, _ := json.Marshal(body)
	path = strings.TrimPrefix(path, "/")

	var resp *http.Response
	var err error
	if body == nil {
		resp, err = u.client.Get(fmt.Sprintf("%s/%s", u.baseUrl, path))
	} else {
		resp, err = u.client.Post(fmt.Sprintf("%s/%s", u.baseUrl, path), "application/json; charset=utf-8", bytes.NewReader(json))
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return resp.StatusCode, nil
}

func (u *Client) Login(username, password string) error {
	data := map[string]string{
		"username": username,
		"password": password,
	}
	statusCode, err := u.do("/api/login", data)

	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"code": statusCode,
	}).Info("login response")

	if statusCode != 200 {
		return errors.New("Controller returned non 200 status code")
	}

	return nil
}

func (u *Client) AuthUser(deviceMac, site, minutes string) error {
	data := map[string]string{
		"cmd":     "authorize-guest",
		"mac":     deviceMac,
		"minutes": minutes,
	}

	statusCode, err := u.do(fmt.Sprintf("/api/s/%s/cmd/stamgr", site), data)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"code": statusCode,
	}).Info("authorize command response")

	if statusCode != 200 {
		return errors.New("Controller returned non 200 status code")
	}
	return nil

}

func (u *Client) Logout() error {
	_, err := u.do(fmt.Sprintf("%s/logout", u.baseUrl), nil)
	if err != nil {
		return err
	}

	return nil
}
