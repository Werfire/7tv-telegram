package main

import (
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Sadzeih/bttv-telegram/pkg/emote"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token: os.Getenv("TOKEN"),
		Poller: &tb.Webhook{
			Listen: os.Getenv("LISTEN_ADDR"),
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: os.Getenv("PUBLIC_URL"),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		var emotes []emote.Emote
		if q.Text == "" {
			return
		}

		// Parse optional leading count: "5 pepe" → maxResults=5, text="pepe".
		text := q.Text
		maxResults := 20
		if parts := strings.SplitN(text, " ", 2); len(parts) == 2 {
			if n, parseErr := strconv.Atoi(parts[0]); parseErr == nil && n > 0 {
				maxResults = int(math.Min(50, float64(n)))
				text = parts[1]
			}
		}

		// Exact search when query is wrapped in double quotes: "pepe".
		exact := len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"'
		if exact {
			text = text[1 : len(text)-1]
		}

		emotes, err := getEmotes(text, exact)
		if err != nil {
			log.Println(err)
			return
		}
		emotes = emotes[:int(math.Min(float64(maxResults), float64(len(emotes))))]

		results := make(tb.Results, len(emotes))
		for i, e := range emotes {
			var result tb.Result
			switch e.Type {
			case "png":
				result = &tb.PhotoResult{
					URL:      e.URL,
					ThumbURL: e.URL,
					Width:    e.Width,
					Height:   e.Height,
				}
			case "gif", "webp":
				result = &tb.GifResult{
					URL:       e.URL,
					ThumbURL:  e.URL,
					ThumbMIME: "image/gif",
					Width:     e.Width,
					Height:    e.Height,
				}
			default:
				result = nil
			}

			if result == nil {
				continue
			}
			result.SetResultID(strconv.Itoa(i + 1))
			results[i] = result
		}

		err = b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 0,
		})
		if err != nil {
			log.Println(err)
		}
	})

	b.Start()
}
