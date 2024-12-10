package suno_api_go

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var cookie = os.Getenv("suno_cookie")

const proxy = "http://127.0.0.1:1080"

// TestBillingInfo 账户信息
func TestBillingInfo(t *testing.T) {
	c := NewClient(Config{
		TimeOut: 10,
		Proxy:   proxy,
		Cookie:  cookie,
	})
	resp, err := c.BillingInfo()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", resp)
}

// TestGenerateByPrompt 根据提示词生成歌曲
func TestGenerateByPrompt(t *testing.T) {
	c := NewClient(Config{
		TimeOut: 10,
		Proxy:   proxy,
		Cookie:  cookie,
	})
	generate, err := c.Generate(GenerateRequest{
		Prompt:           "清风醉",
		Mv:               "chirp-v3-0",
		MakeInstrumental: false,
		Title:            "乡村音乐",
		Tags:             "emotional rap",
	})
	if err != nil {
		return
	}
	var (
		ids         []string
		completeMap = make(map[string]interface{})
		channel     = make(chan struct{})
	)
	for _, v := range generate.Clips {
		ids = append(ids, v.Id)
	}
	tm := time.NewTimer(2 * time.Second)
	go func() {
		defer close(channel)
		for {
			if len(completeMap) == len(ids) {
				break
			}
			task, err := c.GenerateTask(ids)
			if err != nil {
				panic(err)
			}
			if task == nil {
				panic("task is nil")
			}
			for _, v := range task {
				// 表示完成
				if v.Status == Complete {
					completeMap[v.Id] = v
				}
			}
			select {
			case <-tm.C:
				continue
			}
		}
	}()
	<-channel
	fmt.Printf("%+v", completeMap)
}

// TestGenerateByGpt 根据GPT生成歌曲
func TestGenerateByGpt(t *testing.T) {
	c := NewClient(Config{
		TimeOut: 10,
		Proxy:   proxy,
		Cookie:  cookie,
	})
	generate, err := c.Generate(GenerateRequest{
		GptDescriptionPrompt: "乡村音乐",
		Mv:                   "chirp-v3-0",
		MakeInstrumental:     false,
		Title:                "乡村音乐",
		Tags:                 "emotional rap",
	})
	if err != nil {
		return
	}
	var (
		ids         []string
		completeMap = make(map[string]interface{})
		channel     = make(chan struct{})
	)
	for _, v := range generate.Clips {
		ids = append(ids, v.Id)
	}
	tm := time.NewTimer(2 * time.Second)
	go func() {
		defer close(channel)
		for {
			if len(completeMap) == len(ids) {
				break
			}
			task, err := c.GenerateTask(ids)
			if err != nil {
				panic(err)
			}
			if task == nil {
				panic("task is nil")
			}
			for _, v := range task {
				// 表示完成
				if v.Status == Complete {
					completeMap[v.Id] = v
				}
			}
			select {
			case <-tm.C:
				continue
			}
		}
	}()
	<-channel
	fmt.Printf("%+v", completeMap)
}

// TestGenerateLyrics 歌词优化
func TestGenerateLyrics(t *testing.T) {
	c := NewClient(Config{
		TimeOut: 10,
		Proxy:   proxy,
		Cookie:  cookie,
	})
	id, err := c.GenerateLyrics("清风醉")
	if err != nil {
		panic(err)
	}
	if id == "" {
		panic("id is empty")
	}
	var (
		lyricsInfo *GenerateLyricsResponse
		channel    = make(chan struct{})
	)
	go func() {
		defer close(channel)
		for {
			resp, err := c.GetFormatLyrics(id)
			if err != nil {
				return
			}
			if resp.Status == Complete {
				lyricsInfo = resp
				break
			}
		}
	}()
	<-channel
	if lyricsInfo != nil {
		fmt.Printf("title %s,text %s", lyricsInfo.Title, lyricsInfo.Text)
	}
}

// TestGenerateLyricsPair 新版歌词测试
func TestGenerateLyricsPair(t *testing.T) {
	c := NewClient(Config{
		TimeOut: 10,
		Proxy:   proxy,
		Cookie:  cookie,
	})
	result, err := c.GenerateLyricsPair("清风醉")
	if err != nil {
		panic(err)
	}
	if result == nil {
		panic("获取歌词失败")
	}
	var (
		lyricsInfo *GenerateLyricsResponse
		channel    = make(chan struct{})
	)
	go func() {
		defer close(channel)
		for {
			resp, err := c.GetFormatLyrics(result.LyricsAId)
			if err != nil {
				return
			}
			if resp.Status == Complete {
				lyricsInfo = resp
				break
			}
		}
	}()
	<-channel
	if lyricsInfo != nil {
		fmt.Printf("title %s,text %s", lyricsInfo.Title, lyricsInfo.Text)
	}
}

// TestExtendMusic 扩展音乐
func TestExtendMusic(t *testing.T) {
	c := NewClient(Config{
		TimeOut: 10,
		Proxy:   proxy,
		Cookie:  cookie,
	})
	generate, err := c.Generate(GenerateRequest{
		Prompt:         "[Verse]\n在田野里漫步\n看着风儿轻轻吹过\n远方的山峦消失在云雾里\n我心中的宁静永不磨灭\n\n[Verse 2]\n河水流淌在田野间\n闪烁着金色的阳光\n大地为我敞开怀抱\n欢迎我去感受大自然\n\n[Chorus]\n这就是我爱上的田野风光\n大自然的韵律让我心生憧憬\n在这片广袤的土地上\n我找到了属于我的快乐",
		Mv:             "chirp-v3-0",
		Title:          "乡村音乐",
		Tags:           "emotional rap",
		ContinueClipId: "6b704b7f-3629-431e-b1e8-1d532a92d448",
		ContinueAt:     120,
		GenerationType: GenerationTypeText.ToString(),
		Task:           "extend",
	})
	if err != nil {
		return
	}
	var (
		ids         []string
		completeMap = make(map[string]interface{})
		channel     = make(chan struct{})
	)
	for _, v := range generate.Clips {
		ids = append(ids, v.Id)
	}
	tm := time.NewTimer(2 * time.Second)
	go func() {
		defer close(channel)
		for {
			if len(completeMap) == len(ids) {
				break
			}
			task, err := c.GenerateTask(ids)
			if err != nil {
				panic(err)
			}
			if task == nil {
				panic("task is nil")
			}
			for _, v := range task {
				// 表示完成
				if v.Status == Complete {
					completeMap[v.Id] = v
				}
			}
			select {
			case <-tm.C:
				continue
			}
		}
	}()
	<-channel
	fmt.Printf("%+v", completeMap)
}
