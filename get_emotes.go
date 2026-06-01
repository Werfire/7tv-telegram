package main

import (
	"github.com/Sadzeih/bttv-telegram/pkg/emote"
	"github.com/Sadzeih/bttv-telegram/pkg/seventv"
)

func getEmotes(query string, channel string, exact bool) ([]emote.Emote, error) {
	var raw []seventv.Emote
	var err error

	if channel != "" {
		if query == "" {
			return nil, nil
		}
		raw, err = seventv.ChannelEmotes(channel)
		if err != nil {
			return nil, err
		}
		emotes := make([]emote.Emote, len(raw))
		for i, e := range raw {
			emotes[i] = e.Convert()
		}
		if exact {
			return emote.ExactSearchEmotes(query, emotes), nil
		}
		return emote.ChannelSearchEmotes(query, emotes), nil
	}

	raw, err = seventv.SearchEmotes(query)
	if err != nil {
		return nil, err
	}
	emotes := make([]emote.Emote, len(raw))
	for i, e := range raw {
		emotes[i] = e.Convert()
	}
	if exact {
		return emote.ExactSearchEmotes(query, emotes), nil
	}
	return emote.SearchEmotes(query, emotes), nil
}
