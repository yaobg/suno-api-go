# suno-api-go

注：非官方sdk

## Installation
```
# Go Modules
require github.com/yaobg/suno-api-go v0.0.1
```

## Usage

The following samples will assist you to become as comfortable as possible with resty library.

```api
// Import resty into your code and refer it as `suno`.
import "github.com/yaobg/suno-api-go"
```

## cookie

![image-20240407111824834.png](image-20240407111824834.png)

## Example

### 账户信息

```go
func TestBillingInfo(t *testing.T) {
c := NewClient(&Config{
TimeOut: 10,
Proxy:   "127.0.0.1:1080",
Cookie:  cookie,
})
resp, err := c.BillingInfo()
if err != nil {
panic(err)
}
fmt.Printf("%+v", resp)
}
```

### 生成歌曲

```
func TestGenerate(t *testing.T) {
	c := NewClient(&Config{
		TimeOut: 10,
		Proxy:   "127.0.0.1:1080",
		Cookie:  cookie,
	})
	generate, err := c.Generate(&GenerateRequest{
		GptDescriptionPrompt: "乡村音乐",
		Mv:                   "chirp-v3-0",
		MakeInstrumental:     false,
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
			case <-time.After(time.Second * 2):
				continue
			}
		}
	}()
	<-channel
	fmt.Printf("%+v", completeMap)
}
```

### 歌词优化

```
func TestGenerateLyrics(t *testing.T) {
	c := NewClient(&Config{
		TimeOut: 10,
		Proxy:   "127.0.0.1:1080",
		Cookie:  cookie,
	})
	id, err := c.GenerateLyrics("Verse 1:\\n在黑暗的深渊里\\n我心中的火焰燃烧\\n肆意挥洒着狂野的力量\\n毫不畏惧,勇往直前\\n\\nChorus:\\n狂风呼啸,雷电交加\\n潜藏在内心的怒火释放\\n无尽的痛苦,无尽的恐惧\\n唯有摇滚乐撼动心魂")
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
			select {
			case <-time.After(time.Second * 2):
				continue
			case <-time.After(time.Second * 10):
				return
			}
		}
	}()
	<-channel
	if lyricsInfo != nil {
		fmt.Printf("title %s,text %s", lyricsInfo.Title, lyricsInfo.Text)
	}
}
```

