# suno-api-go

注：suno目前没有go的sdk调用，这个是按照官方web请求封装的一个sdk

## Installation

```
# Go Modules
require github.com/yaobg/suno-api-go v0.0.5
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
package main

import (
	"fmt"
	suno "github.com/yaobg/suno-api-go"
)

const cookie = ""

func main() {
	client := suno.NewClient(suno.Config{
		Proxy:  "127.0.0.1:1080",
		Cookie: cookie,
	})
	resp, err := client.BillingInfo()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", resp)
}

```

### 生成歌曲

```go
package main

import (
	"fmt"
	suno "github.com/yaobg/suno-api-go"
	"time"
)

const cookie = ""

func main() {
	client := suno.NewClient(suno.Config{
		Proxy:  "127.0.0.1:1080",
		Cookie: cookie,
	})
	generate, err := client.Generate(suno.GenerateRequest{
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
			task, err := client.GenerateTask(ids)
			if err != nil {
				panic(err)
			}
			if task == nil {
				panic("task is nil")
			}
			for _, v := range task {
				// 表示完成
				if v.Status == suno.Complete {
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

```go
package main

import (
	"fmt"
	suno "github.com/yaobg/suno-api-go"
	"time"
)

const cookie = ""

func main() {
	client := suno.NewClient(suno.Config{
		Proxy:  "127.0.0.1:1080",
		Cookie: cookie,
	})
	id, err := client.GenerateLyrics("Verse 1:\\n在黑暗的深渊里\\n我心中的火焰燃烧\\n肆意挥洒着狂野的力量\\n毫不畏惧,勇往直前\\n\\nChorus:\\n狂风呼啸,雷电交加\\n潜藏在内心的怒火释放\\n无尽的痛苦,无尽的恐惧\\n唯有摇滚乐撼动心魂")
	if err != nil {
		panic(err)
	}
	if id == "" {
		panic("id is empty")
	}
	var (
		lyricsInfo *suno.GenerateLyricsResponse
		channel    = make(chan struct{})
	)
	go func() {
		defer close(channel)
		for {
			resp, err := client.GetFormatLyrics(id)
			if err != nil {
				return
			}
			if resp.Status == suno.Complete {
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
### 歌曲extend
```go
package main

import (
	"fmt"
	suno "github.com/yaobg/suno-api-go"
	"time"
)

const cookie = ""

func main() {
	client := suno.NewClient(suno.Config{
		Proxy:  "127.0.0.1:1080",
		Cookie: cookie,
	})
	generate, err := client.Generate(suno.GenerateRequest{
		Prompt:         "[Verse]\\nLate at night\\nWhen the stars are bright\\nWe hit the streets\\nReady for a wild ride\\nNo destination\\nJust chasing the noise\\nWith the moon as our guide\\nWe let go of our poise\\n\\nDown the boulevard\\nThe city's so alive\\nThe neon lights flicker\\nAs we're taking a dive\\nWe dance to the rhythm\\nLike nobody's watching\\nHeartbeats syncing\\nBodies pulsating\\n\\n[Verse 2]\\nThrough the alleys\\nWe're chasing the thrill\\nFeeding off the energy\\nCan't get our fill\\nEvery corner turned is a brand-new surprise\\nWe're lost in the maze\\nBut we don't realize",
		Mv:             "chirp-v3-0",
		Title:          "乡村音乐",
		Tags:           "emotional rap",
		ContinueClipId: "6384daa1-2725-422d-b489-e9916c9dd68c",
		ContinueAt:     120,
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
			task, err := client.GenerateTask(ids)
			if err != nil {
				panic(err)
			}
			if task == nil {
				panic("task is nil")
			}
			for _, v := range task {
				// 表示完成
				if v.Status == suno.Complete {
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
