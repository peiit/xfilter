package wordlist

import (
	"log"

	jsoniter "github.com/sniperkit/xutil/plugin/format/json"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type HitResponse struct {
	Hit   bool `json:"hit"`
	Level int  `json:"level"`
}

type PredictResponse struct {
	Label int32  `json:"label"`
	Text  string `json:"text"`
}

func RenderJson(res HitResponse) []byte {
	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Printf("Json decode error.")
	}
	return jsonData
}

func DecodeJson(res []byte) PredictResponse {
	var s PredictResponse
	json.Unmarshal(res, &s)
	return s
}
