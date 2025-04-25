package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type TweetData struct {
	Tweet struct {
		ID        string `json:"id"`
		FullText  string `json:"full_text"`
		CreatedAt string `json:"created_at"`
	} `json:"tweet"`
}

func main() {
	startFromID := flag.String("start", "", "このIDから削除を開始します")
	flag.Parse()

	tweetsData, err := os.ReadFile("tweets.js")
	if err != nil {
		log.Fatalf("tweets.jsの読み込みに失敗しました: %v", err)
	}

	savedData, err := os.ReadFile("final_saved_selected_tweets.txt")
	if err != nil {
		log.Fatalf("final_saved_selected_tweets.txtの読み込みに失敗しました: %v", err)
	}

	savedIds := make(map[string]bool)
	for _, id := range strings.Split(string(savedData), "\n") {
		id = strings.TrimSpace(id)
		if id != "" {
			savedIds[id] = true
		}
	}

	var tweets []TweetData
	if err := json.Unmarshal(tweetsData, &tweets); err != nil {
		log.Fatalf("JSONの解析に失敗しました: %v", err)
	}

	var deleteTargets []string
	for _, tweet := range tweets {
		if !savedIds[tweet.Tweet.ID] {
			deleteTargets = append(deleteTargets, tweet.Tweet.ID)
		}
	}

	sort.Slice(deleteTargets, func(i, j int) bool {
		id1 := new(big.Int)
		id2 := new(big.Int)
		id1.SetString(deleteTargets[i], 10)
		id2.SetString(deleteTargets[j], 10)
		return id1.Cmp(id2) > 0
	})

	if *startFromID != "" {
		found := false
		for i, id := range deleteTargets {
			if id == *startFromID {
				deleteTargets = deleteTargets[i:]
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("指定されたID %s が見つかりませんでした", *startFromID)
		}
	}

	commandData, err := os.ReadFile("tweet_delete_command.dat")
	if err != nil {
		log.Fatalf("tweet_delete_command.datの読み込みに失敗しました: %v", err)
	}
	commandTemplate := string(commandData)

	for _, id := range deleteTargets {
		command := strings.ReplaceAll(commandTemplate, "TWEET_ID", id)

		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("ツイートの削除に失敗しました (ID: %s): %v\n出力: %s", id, err, string(output))
			continue
		}

		fmt.Printf("ツイートを削除しました (ID: %s)\n", id)
		
		time.Sleep(500 * time.Millisecond)
	}
} 