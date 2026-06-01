package seventv

import (
	"context"
	"fmt"
	graphql "github.com/hasura/go-graphql-client"
	"net/http"
	"strings"
)

// ObjectID is a custom scalar type matching 7TV's ObjectID GQL type.
type ObjectID string

// ChannelEmotes fetches all emotes from a channel's emote sets by 7TV username (case-insensitive).
func ChannelEmotes(username string) ([]Emote, error) {
	client := graphql.NewClient(gqlEndpoint, http.DefaultClient)
	lower := strings.ToLower(username)

	// Step 1: search for users and find an exact username match.
	usersQ := struct {
		Users []struct {
			ID       string
			Username string
		} `graphql:"users(query: $query)"`
	}{}
	if err := client.Query(context.Background(), &usersQ, map[string]interface{}{
		"query": graphql.String(lower),
	}); err != nil {
		return nil, fmt.Errorf("could not search 7TV users for %q: %w", username, err)
	}

	var userID string
	for _, u := range usersQ.Users {
		if strings.ToLower(u.Username) == lower {
			userID = u.ID
			break
		}
	}
	if userID == "" {
		return nil, fmt.Errorf("channel %q not found on 7TV", username)
	}

	// Step 2: fetch emote sets for the resolved user ID.
	emoteSetQ := struct {
		User struct {
			EmoteSets []struct {
				Emotes []struct {
					Name string
					Data struct {
						Animated bool
						Host     struct {
							Url   string
							Files []struct {
								Width  int
								Height int
							}
						}
					}
				}
			} `graphql:"emote_sets"`
		} `graphql:"user(id: $id)"`
	}{}
	if err := client.Query(context.Background(), &emoteSetQ, map[string]interface{}{
		"id": ObjectID(userID),
	}); err != nil {
		return nil, fmt.Errorf("could not fetch emotes for channel %q: %w", username, err)
	}

	seen := make(map[string]bool)
	var emotes []Emote
	for _, set := range emoteSetQ.User.EmoteSets {
		for _, e := range set.Emotes {
			if seen[e.Name] || e.Data.Host.Url == "" || len(e.Data.Host.Files) == 0 {
				continue
			}
			seen[e.Name] = true
			em := Emote{Name: e.Name}
			file := e.Data.Host.Files[len(e.Data.Host.Files)-1]
			em.Width = file.Width
			em.Height = file.Height
			em.Mime = "png"
			if e.Data.Animated {
				em.Mime = "gif"
			}
			em.URL = fmt.Sprintf("https:%s/%s.%s", e.Data.Host.Url, "4x", em.Mime)
			emotes = append(emotes, em)
		}
	}
	return emotes, nil
}
