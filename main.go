package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"
)

type TweetData struct {
	Tweet struct {
		ID        string `json:"id"`
		FullText  string `json:"full_text"`
		CreatedAt string `json:"created_at"`
	} `json:"tweet"`
}

type TemplateData struct {
	Tweets        []TweetData
	SelectedTweets string
}

func main() {
	data, err := os.ReadFile("tweets.js")
	if err != nil {
		log.Fatalf("tweets.jsの読み込みに失敗しました: %v", err)
	}

	var tweets []TweetData
	if err := json.Unmarshal(data, &tweets); err != nil {
		log.Fatalf("JSONの解析に失敗しました: %v", err)
	}

	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].Tweet.ID < tweets[j].Tweet.ID
	})

	var selectedTweets string
	if savedData, err := os.ReadFile("saved_selected_tweets.txt"); err == nil {
		selectedTweets = string(savedData)
	}

	templateData := TemplateData{
		Tweets:        tweets,
		SelectedTweets: selectedTweets,
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>ツイート選択</title>
    <style>
        .tweet {
            border: 1px solid #ccc;
            margin: 10px;
            padding: 10px;
            cursor: pointer;
        }
        .tweet.selected {
            background-color: #e6f7ff;
        }
        .tweet-text {
            white-space: pre-wrap;
        }
        .save-button {
            margin: 20px;
            padding: 10px 20px;
            font-size: 16px;
            cursor: pointer;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 4px;
            position: fixed;
            bottom: 20px;
            right: 20px;
            z-index: 1000;
        }
        .save-button:hover {
            background-color: #45a049;
        }
        body {
            padding-bottom: 80px;
        }
    </style>
</head>
<body>
    <h1>ツイート選択</h1>
    <div id="tweets">
        {{range .Tweets}}
        <div class="tweet" id="{{.Tweet.ID}}">
            <p class="tweet-text">{{.Tweet.FullText}}</p>
            <small>{{.Tweet.CreatedAt}}</small>
        </div>
        {{end}}
    </div>
    <button class="save-button" id="saveButton">保存する</button>
    <script>
        let selectedTweets = new Set();
        
        const savedIds = ` + "`{{.SelectedTweets}}`" + `.split('\n').filter(id => id.trim() !== '');
        savedIds.forEach(id => {
            selectedTweets.add(id);
            const tweet = document.querySelector('div[id=\'' + id + '\']');
            if (tweet) {
                tweet.classList.add('selected');
            }
        });
        
        document.querySelectorAll('.tweet').forEach(tweet => {
            tweet.addEventListener('click', () => {
                const id = tweet.id;
                if (selectedTweets.has(id)) {
                    selectedTweets.delete(id);
                    tweet.classList.remove('selected');
                } else {
                    selectedTweets.add(id);
                    tweet.classList.add('selected');
                }
            });
        });

        document.getElementById('saveButton').addEventListener('click', () => {
            if (selectedTweets.size === 0) {
                alert('選択されたツイートがありません');
                return;
            }
            
            const ids = Array.from(selectedTweets).join('\n');
            const blob = new Blob([ids], { type: 'text/plain' });
            const a = document.createElement('a');
            a.href = URL.createObjectURL(blob);
            a.download = 'saved_selected_tweets.txt';
            a.click();
            URL.revokeObjectURL(a.href);
        });
    </script>
</body>
</html>`

	t, err := template.New("tweets").Parse(tmpl)
	if err != nil {
		log.Fatalf("テンプレートの解析に失敗しました: %v", err)
	}

	output, err := os.Create("tweets.html")
	if err != nil {
		log.Fatalf("HTMLファイルの作成に失敗しました: %v", err)
	}
	defer output.Close()

	if err := t.Execute(output, templateData); err != nil {
		log.Fatalf("テンプレートの実行に失敗しました: %v", err)
	}

	fmt.Println("tweets.htmlを生成しました")
} 