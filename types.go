package suno_api_go

import (
	"github.com/go-resty/resty/v2"
	"time"
)

// Client client
type Client struct {
	Config
	client *resty.Client
}

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
	History              []History   `json:"history"`
	ConcatHistory        interface{} `json:"concat_history"`
	Type                 string      `json:"type"`
	Duration             float64     `json:"duration"`
	RefundCredits        bool        `json:"refund_credits"`
	Stream               bool        `json:"stream"`
	ErrorType            interface{} `json:"error_type"`
	ErrorMessage         interface{} `json:"error_message"`
}

type History struct {
	Id         string `json:"id"`
	ContinueAt int    `json:"continue_at"`
}

// Config 配置
type Config struct {
	Proxy       string
	Cookie      string
	TimeOut     int64
	GenerateUrl string // 生成歌词Url地址
}

type GenerationType string

const (
	GenerationTypeText    GenerationType = "TEXT"
	GenerationTypeAUDIO   GenerationType = "AUDIO"
	GenerationTypeIMAGE   GenerationType = "IMAGE"
	GenerationTypeVIDEO   GenerationType = "VIDEO"
	GenerationTypeTWITTER GenerationType = "TWITTER"
)

func (g GenerationType) ToString() string {
	return string(g)
}

// GenerateRequest GenerateRequest
type GenerateRequest struct {
	GptDescriptionPrompt string `json:"gpt_description_prompt,omitempty"` //gpt提示词
	Mv                   string `json:"mv"`                               //版本
	Prompt               string `json:"prompt"`                           //提示词
	MakeInstrumental     bool   `json:"make_instrumental"`                //是否只要音乐
	Title                string `json:"title"`                            //标题
	Tags                 string `json:"tags"`                             //风格
	ContinueAt           int    `json:"continue_at"`                      //扩展歌词对接时间
	ContinueClipId       string `json:"continue_clip_id"`                 //扩展歌词id
	Task                 string `json:"task"`                             //任务类型 扩展：extend
	GenerationType       string `json:"generation_type"`                  //生成类型 文本：text 音频：AUDIO 图片：IMAGE 视频：VIDEO
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
	Detail interface{} `json:"detail"`
}
type GenerateLyricsResponse struct {
	Text   string `json:"text"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// BillingInfoResponse 账单信息
type BillingInfoResponse struct {
	IsActive         bool  `json:"is_active"`
	IsPastDue        bool  `json:"is_past_due"`
	Credits          int64 `json:"credits"`
	SubscriptionType bool  `json:"subscription_type"`
	TotalCreditsLeft int64 `json:"total_credits_left"`
}

type Sessions struct {
	Object                   string          `json:"object"`
	ID                       string          `json:"id"`
	Status                   string          `json:"status"`
	ExpireAt                 int64           `json:"expire_at"`
	AbandonAt                int64           `json:"abandon_at"`
	LastActiveAt             int64           `json:"last_active_at"`
	LastActiveOrganizationID *string         `json:"last_active_organization_id"`
	Actor                    *string         `json:"actor"`
	FactorVerificationAge    []int           `json:"factor_verification_age"`
	CreatedAt                int64           `json:"created_at"`
	UpdatedAt                int64           `json:"updated_at"`
	LastActiveToken          LastActiveToken `json:"last_active_token"`
}
type LastActiveToken struct {
	Object string `json:"object"`
	Jwt    string `json:"jwt"`
}

// LyricsPairResponse 新版歌词响应值
type LyricsPairResponse struct {
	LyricsAId       string `json:"lyrics_a_id"`
	LyricsBId       string `json:"lyrics_b_id"`
	LyricsRequestId string `json:"lyrics_request_id"`
}

// TokenResponse token
type TokenResponse struct {
	Response struct {
		Object              string      `json:"object"`
		ID                  string      `json:"id"`
		Sessions            []Sessions  `json:"sessions"`
		SignIn              interface{} `json:"sign_in"`
		SignUp              interface{} `json:"sign_up"`
		LastActiveSessionID string      `json:"last_active_session_id"`
		CookieExpiresAt     *string     `json:"cookie_expires_at"`
		CreatedAt           int64       `json:"created_at"`
		UpdatedAt           int64       `json:"updated_at"`
	} `json:"response"`
	Client interface{} `json:"client"`
}
