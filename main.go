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
		if q.Text == "" {
			return
		}

		// Parse page offset sent by Telegram on scroll.
		pageOffset := 0
		if q.Offset != "" {
			if n, parseErr := strconv.Atoi(q.Offset); parseErr == nil && n >= 0 {
				pageOffset = n
			}
		}

		text := q.Text
		maxResults := 20
		countSet := false

		// Parse optional leading count: "5 pepe" → maxResults=5, text="pepe".
		if parts := strings.SplitN(text, " ", 2); len(parts) == 2 {
			if n, parseErr := strconv.Atoi(parts[0]); parseErr == nil && n > 0 {
				maxResults = int(math.Min(50, float64(n)))
				text = parts[1]
				countSet = true
			}
		}

		// Parse optional channel prefix: "@xqc pepe" → channel="xqc", text="pepe".
		channel := ""
		if strings.HasPrefix(text, "@") {
			parts := strings.SplitN(text, " ", 2)
			channel = parts[0][1:]
			if len(parts) == 2 {
				text = parts[1]
			} else {
				text = ""
			}
		}

		// Parse * prefix: show emote name as title in results (ignored for exact search).
		showNames := false
		if strings.HasPrefix(text, "*") {
			showNames = true
			text = text[1:]
		}

		// Exact search when query is wrapped in double quotes: "pepe".
		exact := len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"'
		if exact {
			text = text[1 : len(text)-1]
			showNames = false
		}

		// Global search requires a non-empty query.
		if text == "" && channel == "" {
			return
		}

		var emotes []emote.Emote
		emotes, err = getEmotes(text, channel, exact)

		// On subsequent pages, silently stop on error or exhaustion.
		if pageOffset > 0 && (err != nil || len(emotes) == 0) {
			b.Answer(q, &tb.QueryResponse{CacheTime: 0})
			return
		}

		// On the first page, show error/empty states as article results.
		if err != nil || len(emotes) == 0 {
			var article *tb.ArticleResult
			if err != nil {
				log.Println(err)
				article = &tb.ArticleResult{
					Title: "❌ Channel not found on 7TV",
					Text:  "❌ Channel not found on 7TV",
				}
			} else {
				article = &tb.ArticleResult{
					Title: "🔍 No emotes found",
					Text:  "🔍 No emotes found",
				}
			}
			article.SetResultID("1")
			if answerErr := b.Answer(q, &tb.QueryResponse{Results: tb.Results{article}, CacheTime: 0}); answerErr != nil {
				log.Println(answerErr)
			}
			return
		}

		// Channel results: paginate. Global results: apply limit directly.
		var nextOffset string
		if channel != "" {
			pageSize := maxResults
			if !countSet {
				pageSize = 50
			}
			total := len(emotes)
			start := pageOffset
			if start >= total {
				b.Answer(q, &tb.QueryResponse{CacheTime: 0})
				return
			}
			end := int(math.Min(float64(start+pageSize), float64(total)))
			emotes = emotes[start:end]
			if end < total {
				nextOffset = strconv.Itoa(end)
			}
		} else {
			emotes = emotes[:int(math.Min(float64(maxResults), float64(len(emotes))))]
		}

		results := make(tb.Results, len(emotes))
		for i, e := range emotes {
			var result tb.Result
			switch e.Type {
			case "png":
				r := &tb.PhotoResult{
					URL:      e.URL,
					ThumbURL: e.URL,
					Width:    e.Width,
					Height:   e.Height,
				}
				if showNames {
					r.Title = e.Name
				}
				result = r
			case "gif", "webp":
				r := &tb.GifResult{
					URL:       e.URL,
					ThumbURL:  e.URL,
					ThumbMIME: "image/gif",
					Width:     e.Width,
					Height:    e.Height,
				}
				if showNames {
					r.Title = e.Name
				}
				result = r
			default:
				result = nil
			}

			if result == nil {
				continue
			}
			result.SetResultID(strconv.Itoa(pageOffset + i + 1))
			results[i] = result
		}

		if answerErr := b.Answer(q, &tb.QueryResponse{
			Results:    results,
			NextOffset: nextOffset,
			CacheTime:  0,
		}); answerErr != nil {
			log.Println(answerErr)
		}
	})

	b.Start()
}
