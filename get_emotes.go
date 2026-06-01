package main

import (
	"github.com/Sadzeih/bttv-telegram/pkg/emote"
	"github.com/Sadzeih/bttv-telegram/pkg/seventv"
)

func getEmotes(query string, exact bool) ([]emote.Emote, error) {
	stv, err := seventv.SearchEmotes(query)
	if err != nil {
		return nil, err
	}

	emotes := make([]emote.Emote, len(stv))
	for i, e := range stv {
		emotes[i] = e.Convert()
	}

	if exact {
		return emote.ExactSearchEmotes(query, emotes), nil
	}
	return emote.SearchEmotes(query, emotes), nil
}
