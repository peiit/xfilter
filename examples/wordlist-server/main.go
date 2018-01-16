package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sniperkit/xfilter/backend/wordlist"
)

func main() {
	// load dict and watch file change
	wordlist.InitLoader()

	log.Println("Run word_filter app on", 8001)
	// routes
	http.HandleFunc("/word/is_valid", wordlist.VerifyWordsHandler)
	// run server
	http.ListenAndServe(":8001", nil)
}

// handle funcs
func VerifyWordsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		text := req.FormValue("v")
		var hit bool
		var review_level int

		level := VerifyWords(text)
		if level == 0 {
			hit = false
			review_level = 0
		} else {
			hit = true
			// if level equals to 1, need to predict text label
			if level == 1 {
				label := PredictText(text)
				if label == 0 {
					review_level = 2 // need ban
				} else {
					review_level = level
				}
			} else {
				review_level = level
			}
		}
		log.Printf("[INFO] text: " + text + " hit: " + fmt.Sprintf("%t", hit) + " review_level: " + fmt.Sprintf("%d", review_level))
		res := HitResponse{Hit: hit, Level: review_level}
		w.Write(RenderJson(res))
	}
}
