# Tweet Cleaner

A program to delete tweets from your Twitter archive. This program reads tweets from a Twitter archive file and deletes.

## Preparation

- Twitter archive files:
  - `tweets.js`: Your Twitter archive file
  - `tweet_delete_command.dat`: File containing the tweet deletion command template

## Usage

### 1 : main.go

1. Place `tweets.js` at same directory and run `main.go`:
```bash
go run main.go
```

2. Reads tweets from `tweets.js`.

3. Generates an HTML file (`tweets.html`) for selecting tweets to keep.

4. Open `tweets.html` and select tweets you want to keep.

5. Saves the selected tweet IDs to `saved_selected_tweets.txt`.

### 2 : delete_tweets.go

1. Place `final_saved_selected_tweets.txt` and `tweets.js` and run `delete_tweets.go`:
   ```bash
   go run delete_tweets.go
   ```
2. After that tweets except selected tweet IDs will be deleted.

3. To resume deletion from a specific tweet ID:
   ```bash
   go run delete_tweets.go -start TWEET_ID
   ```
