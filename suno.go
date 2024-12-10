package suno_api_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"strings"
	"time"
)

const (
	Queued    = "queued"
	Streaming = "streaming"
	Complete  = "complete"
)

const baseUrl = "https://studio-api.suno.ai"

type option func(config *Config)

// NewClient client
func NewClient(c Config) *Client {
	var (
		timeOut time.Duration
	)
	if c.TimeOut == 0 {
		timeOut = 10 * time.Second
	}
	client := resty.New().SetTimeout(timeOut)
	if c.Proxy != "" {
		client.SetProxy(c.Proxy)
	}
	return &Client{
		Config: c,
		client: client,
	}
}

func (g generateError) ToString() string {
	data, _ := json.Marshal(g)
	return string(data)
}

// Generate Generate
func (s *Client) Generate(req GenerateRequest) (data *GenerateResponse, err error) {
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	path := baseUrl + "/api/generate/v2/"
	if s.GenerateUrl != "" {
		path = s.GenerateUrl
	}
	var (
		result    GenerateResponse
		resultErr generateError
	)
	payload, _ := json.Marshal(req)
	fmt.Println(string(payload))
	// 需要忽略错误，因为扩展音乐的时候,http 200的时候,err不为空
	r, _ := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetBody(payload).
		SetResult(&result).
		SetError(resultErr).
		Post(path)
	fmt.Println(string(r.Body()))
	if r.StatusCode() != 200 {
		return nil, errors.New(r.String())
	}
	if resultErr.Detail != nil {
		return nil, errors.New(resultErr.ToString())
	}
	return &result, nil
}

// GenerateTask 获取任务
// 需要自己建立监听
func (s *Client) GenerateTask(ids []string) (data []Clips, err error) {
	var (
		result []Clips
	)
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	_, err = s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", s.Cookie).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetResult(&result).
		Get(fmt.Sprintf("%s/api/feed/?ids=%s", baseUrl, url.PathEscape(strings.Join(ids, ","))))
	return result, nil
}

// GenerateLyrics generate lyrics
func (s *Client) GenerateLyrics(prompt string) (id string, err error) {
	token, err := s.getToken()
	if err != nil {
		return "", err
	}
	path := baseUrl + "/api/generate/lyrics/"
	var req struct {
		Prompt string `json:"prompt"`
	}
	req.Prompt = prompt
	var (
		result = make(map[string]string)
	)
	r, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetBody(req).
		SetResult(&result).
		Post(path)
	if err != nil {
		return "", err
	}
	if r.StatusCode() != 200 {
		return "", errors.New(r.String())
	}
	return result["id"], nil
}

// GenerateLyricsPair generate lyrics-pair
func (s *Client) GenerateLyricsPair(prompt string) (result *LyricsPairResponse, err error) {
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	path := baseUrl + "/api/generate/lyrics-pair/"
	var req struct {
		Prompt      string `json:"prompt"`
		LyricsModel string `json:"lyrics_model"`
	}
	req.Prompt = prompt
	r, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetHeader("User-Agent", "PostmanRuntime/7.42.0").
		SetBody(req).
		SetResult(result).
		Post(path)
	if err != nil {
		return nil, err
	}
	if r.StatusCode() != 200 {
		return nil, errors.New(r.String())
	}
	return result, nil
}

// GetFormatLyrics get Format lyrics
func (s *Client) GetFormatLyrics(id string) (data *GenerateLyricsResponse, err error) {
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	var (
		result GenerateLyricsResponse
	)
	path := fmt.Sprintf("%s/api/generate/lyrics/%s", baseUrl, id)
	r, err := s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetResult(&result).
		Get(path)
	if err != nil {
		return nil, err
	}
	if r.StatusCode() != 200 {
		return nil, errors.New(r.String())
	}
	return &result, nil
}

// getToken get token
func (s *Client) getToken() (string, error) {
	if s.Cookie == "" {
		return "", errors.New("cookie is empty")
	}
	path := "https://clerk.suno.com/v1/client?_clerk_js_version=4.72.2"
	var (
		tokenRes TokenResponse
	)
	r, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", s.Cookie).
		SetResult(&tokenRes).
		Get(path)
	if err != nil {
		return "", fmt.Errorf("get token Authorization err:%s", err)
	}
	if r.StatusCode() != 200 {
		return "", errors.New(r.String())
	}
	fmt.Printf(string(r.Body()))
	if tokenRes.Response.Sessions == nil || len(tokenRes.Response.Sessions) == 0 {
		return "", fmt.Errorf("get token Authorization")
	}
	token := tokenRes.Response.Sessions[0].LastActiveToken.Jwt
	return token, nil
}

// BillingInfo 获取账单信息
func (s *Client) BillingInfo() (data *BillingInfoResponse, err error) {
	// 根据seesion_id获取token
	path := "https://studio-api.suno.ai/api/billing/info/"
	token, err := s.getToken()
	data = &BillingInfoResponse{}
	r, err := s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetResult(data).
		Get(path)
	if err != nil {
		return nil, err
	}
	if r.StatusCode() != 200 {
		return nil, errors.New(r.String())
	}
	return
}
