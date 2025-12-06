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

func NewClient(url string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	return &Client{
		client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		baseUrl: strings.TrimSuffix(url, "/"),
	}, nil
}

func (u *Client) doGet(path string) (int, error) {
	url := fmt.Sprintf("%s/%s", u.baseUrl, strings.TrimPrefix(path, "/"))
	resp, err := u.client.Get(url)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (u *Client) doPost(path string, body any) (int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to marshal body: %w", err)
	}
	url := fmt.Sprintf("%s/%s", u.baseUrl, strings.TrimPrefix(path, "/"))
	resp, err := u.client.Post(url, "application/json; charset=utf-8", bytes.NewReader(data))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (u *Client) Login(username, password string) error {
	data := map[string]string{
		"username": username,
		"password": password,
	}
	statusCode, err := u.doPost("/api/login", data)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"code": statusCode,
	}).Info("login response")

	if statusCode != 200 {
		return errors.New("controller returned non 200 status code")
	}
	return nil
}

func (u *Client) AuthUser(deviceMac, site, minutes string) error {
	data := map[string]string{
		"cmd":     "authorize-guest",
		"mac":     deviceMac,
		"minutes": minutes,
	}
	statusCode, err := u.doPost(fmt.Sprintf("/api/s/%s/cmd/stamgr", site), data)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"code": statusCode,
	}).Info("authorize command response")

	if statusCode != 200 {
		return errors.New("controller returned non 200 status code")
	}
	return nil
}

func (u *Client) Logout() error {
	_, err := u.doGet("/api/logout")
	return err
}
