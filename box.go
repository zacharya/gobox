package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Config struct {
	Client ClientConfig `json:"client"`
	Events EventsConfig `json:"events"`
}

type ClientConfig struct {
	Token           string          `json:"token"`
	UrlBase         string          `json:"url_base"`
	ClientID        string          `json:"client_id"`
	ClientSecret    string          `json:"client_secret"`
	JWTCustomClaims JWTCustomClaims `json:"jwt_custom_claims"`
}

type EventsConfig struct {
	StartTime  string `json:"start_time"`
	EventLimit string `json:"event_limit"`
}

type JWTCustomClaims struct {
	Iss     string `json:"iss"`
	Sub     string `json:"sub"`
	SubType string `json:"box_sub_type"`
	Aud     string `json:"aud"`
	Jti     string `json:"jti"`
	Exp     int64  `json:"exp"`
	KeyID   string `json:"key_id"`
	jwt.StandardClaims
}

type Client struct {
	HttpClient *http.Client
	BaseUrl    *url.URL
	Token      string
}

func NewClient(tok string, urlBase string) *Client {
	baseURL, _ := url.Parse(urlBase)
	return &Client{
		HttpClient: &http.Client{},
		BaseUrl:    baseURL,
		Token:      tok,
	}
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	_, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	resolvedUrl, err := url.Parse(c.BaseUrl.String() + urlStr)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if body != nil {
		if err = json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, resolvedUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, respStr interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		//fmt.Printf("box client.Do error: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if resp.StatusCode == 429 {
			retrySecs, _ := strconv.Atoi(resp.Header.Get("Retry-after"))
			time.Sleep(time.Duration(retrySecs) * time.Second)
			return nil, fmt.Errorf("rate limited for %d seconds", retrySecs)
		} else {
			return nil, fmt.Errorf("http request failed, resp: %#v", resp)
		}
	}

	if respStr != nil {
		err = json.NewDecoder(resp.Body).Decode(respStr)
	}
	return resp, err
}

func (c *Client) DoWithRetries(req *http.Request, respStr interface{}, retries int) (*http.Response, error) {
	var resp http.Response
	for i := 0; i < retries; i++ {
		resp, err := c.Do(req, respStr)
		if err != nil {
			if i >= retries {
				return nil, fmt.Errorf("http request failed, resp: %#v", resp)
			}
			continue
		}
	}
	return &resp, nil
}

func (c *Client) EventService() *EventService {
	return &EventService{
		Client: c,
	}
}

func (c *Client) FileService() *FileService {
	return &FileService{
		Client: c,
	}
}
