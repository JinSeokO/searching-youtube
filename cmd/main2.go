package main

import (
	"context"
	"flag"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"strings"
)

var (
	query        = flag.String("query", "Google", "Search term")
	developerKey = flag.String("key", "", "developer-key")
	maxResults   = flag.Int64("max-results", 25, "Max YouTube results")
)

func main() {
	flag.Parse()
	if len(strings.TrimSpace(*developerKey)) == 0 {
		log.Fatalf("'key' flag is empty. it is required field")
	}
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(*developerKey))
	if err != nil {
		log.Fatalf("fail to open youtube service : %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List([]string{"id", "snippet"}).
		Q(*query).
		Order("viewCount").
		MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	printIDs("Videos", videos)
	log.Println("===============")
	printIDs("Channels", channels)
	log.Println("===============")
	printIDs("Playlists", playlists)
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	log.Printf("%v:\n", sectionName)
	for id, title := range matches {
		log.Printf("[%v] %v\n", id, title)
	}
	log.Printf("\n\n")
}
