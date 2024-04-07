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

// Clips clips
type Clips struct {
	Detail            string      `json:"detail"`
	Id                string      `json:"id"`
	VideoUrl          string      `json:"video_url"`
	AudioUrl          string      `json:"audio_url"`
	ImageUrl          string      `json:"image_url"`
	ImageLargeUrl     string      `json:"image_large_url"`
	MajorModelVersion string      `json:"major_model_version"`
	ModelName         string      `json:"model_name"`
	Metadata          *Metadata   `json:"metadata"`
	IsLiked           bool        `json:"is_liked"`
	UserId            string      `json:"user_id"`
	IsTrashed         bool        `json:"is_trashed"`
	Reaction          interface{} `json:"reaction"`
	CreatedAt         time.Time   `json:"created_at"`
	Status            string      `json:"status"`
	Title             string      `json:"title"`
	PlayCount         int         `json:"play_count"`
	UpvoteCount       int         `json:"upvote_count"`
	IsPublic          bool        `json:"is_public"`
}

// Metadata Metadata
type Metadata struct {
	Tags                 string      `json:"tags"`
	Prompt               string      `json:"prompt"`
	GptDescriptionPrompt string      `json:"gpt_description_prompt"`
	AudioPromptId        interface{} `json:"audio_prompt_id"`
	History              interface{} `json:"history"`
	ConcatHistory        interface{} `json:"concat_history"`
	Type                 string      `json:"type"`
	Duration             float64     `json:"duration"`
	RefundCredits        bool        `json:"refund_credits"`
	Stream               bool        `json:"stream"`
	ErrorType            interface{} `json:"error_type"`
	ErrorMessage         interface{} `json:"error_message"`
}

// Config 配置
type Config struct {
	Proxy       string
	Cookie      string
	TimeOut     int64
	GenerateUrl string // 生成歌词Url地址
}

type option func(config *Config)

// Client client
type Client struct {
	Config
	client *resty.Client
}

// NewClient client
func NewClient(c *Config) *Client {
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
		Config: *c,
		client: client,
	}
}

// GenerateRequest GenerateRequest
type GenerateRequest struct {
	GptDescriptionPrompt string `json:"gpt_description_prompt"` //gpt提示词
	Mv                   string `json:"mv"`                     //版本
	Prompt               string `json:"prompt"`                 //提示词
	MakeInstrumental     bool   `json:"make_instrumental""`     //是否只要音乐
}

// GenerateResponse	GenerateResponse
type GenerateResponse struct {
	Id                string  `json:"id"`
	Clips             []Clips `json:"clips"`
	Metadata          `json:"metadata"`
	MajorModelVersion string `json:"major_model_version"`
	Status            string `json:"status"`
	CreatedAt         string `json:"created_at"`
	BatchSize         int    `json:"batch_size"`
}
type generateError struct {
	Detail string `json:"detail"`
}

// Generate Generate
func (s *Client) Generate(req *GenerateRequest) (data *GenerateResponse, err error) {
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	url := baseUrl + "/api/generate/v2/"
	if s.GenerateUrl != "" {
		url = s.GenerateUrl
	}
	var (
		result    GenerateResponse
		resultErr generateError
	)
	payload, _ := json.Marshal(req)
	_, err = s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetBody(payload).
		SetResult(&result).
		SetError(resultErr).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resultErr.Detail != "" {
		return nil, errors.New(resultErr.Detail)
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

type GenerateLyricsResponse struct {
	Text   string `json:"text"`
	Title  string `json:"title"`
	Status string `json:"status"`
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

// TokenResponse token 响应值
type TokenResponse struct {
	Response map[string]interface{} `json:"response"`
}

type Sessions struct {
	Id string `json:"id"`
}

// getToken get token
func (s *Client) getToken() (string, error) {
	path := "https://clerk.suno.ai/v1/client"
	var (
		result TokenResponse
	)
	r, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", s.Cookie).
		SetResult(&result).
		Get(path)

	if err != nil {
		return "", err
	}
	if r.StatusCode() != 200 {
		return "", errors.New(r.String())
	}
	if result.Response == nil {
		return "", errors.New("response is empty")
	}
	lastActiveSessionId := result.Response["last_active_session_id"]
	if lastActiveSessionId == "" {
		return "", errors.New("active seesion_id is empty")
	}
	sessions := result.Response["sessions"].([]interface{})
	var token struct {
		JWT string `json:"jwt"`
	}
	session := sessions[0].(map[string]interface{})
	// 根据seesion_id获取token
	path = fmt.Sprintf("https://clerk.suno.ai/v1/client/sessions/%s/tokens", session["id"])
	r, err = s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", s.Cookie).
		SetResult(&token).
		Post(path)
	if err != nil {
		return "", err
	}
	if r.StatusCode() != 200 {
		return "", errors.New(r.String())
	}
	return token.JWT, nil
}

// BillingInfoResponse 账单信息
type BillingInfoResponse struct {
	IsActive         bool  `json:"is_active"`
	IsPastDue        bool  `json:"is_past_due"`
	Credits          int64 `json:"credits"`
	SubscriptionType bool  `json:"subscription_type"`
	TotalCreditsLeft int64 `json:"total_credits_left"`
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
