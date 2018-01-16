package wordlist

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	PREDICT_HOST = "http://127.0.0.1:8006/predict"
)

func PredictText(text string) int32 {
	resp, err := http.PostForm(PREDICT_HOST,
		url.Values{"text": {text}})
	if err != nil {
		// handle error
		log.Printf("[ERROR]: request for Predict Host error!")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Printf("[ERROR]: response from Predict Host error!")
	}
	d := DecodeJson(body)
	log.Printf("[PREDICT-INFO] text: " + d.Text + " label: " + fmt.Sprintf("%d", d.Label))
	return d.Label
}
