package main

import (
	"github.com/Sadzeih/bttv-telegram/pkg/emote"
	"github.com/Sadzeih/bttv-telegram/pkg/seventv"
)

func getEmotes(query string, channel string, exact bool) ([]emote.Emote, error) {
	var raw []seventv.Emote
	var err error

	if channel != "" {
		raw, err = seventv.ChannelEmotes(channel)
	} else {
		raw, err = seventv.SearchEmotes(query)
	}
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
	if query == "" {
		return emotes, nil
	}
	return emote.SearchEmotes(query, emotes), nil
}
