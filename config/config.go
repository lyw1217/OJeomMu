package config

import (
	"log"
	"os"
)

const (
	keyPath   string = "./config/keys.json"
	ServerCrt string = "/etc/letsencrypt/live/mumeog.site/fullchain.pem"
	ServerKey string = "/etc/letsencrypt/live/mumeog.site/privkey.pem"
)

type Kakao_t struct {
	Rest  string `json:"rest_api"`
	JS    string `json:"javascript"`
	Admin string `json:"admin"`
}

type Keys_t struct {
	Kakao Kakao_t `json:"kakao"`
	Newyo struct {
		Apikey string `json:"apikey"`
	} `json:"newyo"`
}

// Load keys from json file
func LoadKeysConfig() Keys_t {
	var k Keys_t

	k.Kakao.Rest = os.Getenv("KAKAO_REST")
	k.Kakao.JS = os.Getenv("KAKAO_JS")
	k.Kakao.Admin = os.Getenv("KAKAO_ADMIN")
	if k.Kakao.Rest == "" {
		log.Fatal("$KAKAO_REST must be set")
	}

	k.Newyo.Apikey = os.Getenv("NEWYO_KEY")
	if k.Newyo.Apikey == "" {
		log.Fatal("$NEWYO_KEY must be set")
	}

	log.Print("Successful loading of Key Info ........")

	return k
}

var Keys Keys_t

func init() {
	Keys = LoadKeysConfig()
}
