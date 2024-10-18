package openai

import (
	"encoding/json"
	"github.com/chinaboard/unifi-captive-portal/pkg/options"
	"net/http"
	"strings"
)

type ChatCompletionRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream      bool    `json:"stream,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type Client struct {
	opt *options.OpenAiOptions
}

func NewClient(opt *options.OpenAiOptions) *Client {
	return &Client{
		opt: opt,
	}
}

func (c *Client) ChatHandler(w http.ResponseWriter, r *http.Request) {
	if c.opt.ApiKey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		return
	}

	var req ChatCompletionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Stream = true
	req.Model = c.opt.Model
	req.Temperature = c.opt.Temperature

	body, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	openaiReq, err := http.NewRequest("POST", strings.TrimSuffix(c.opt.Domain, "/")+"/v1/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	openaiReq.Header.Set("Authorization", "Bearer "+c.opt.ApiKey)
	openaiReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(openaiReq)
	if err != nil {
		http.Error(w, "Failed to get response", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, e := resp.Body.Read(buf)
		if e != nil {
			break
		}
		if n > 0 {
			w.Write(buf[:n])
			flusher.Flush()
		}
	}
}
