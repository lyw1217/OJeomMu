package config

import (
	"encoding/json" // https://pkg.go.dev/encoding/json
	"log"
	"os"
	"path/filepath"
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

type Naver_t struct {
	Id        string `json:"id"`
	Secret    string `json:"secret"`
	NcpId     string `json:"ncp_id"`
	NcpSecret string `json:"ncp_secret"`
}

type Keys_t struct {
	Kakao Kakao_t `json:"kakao"`
	Naver Naver_t `json:"naver"`
	Newyo struct {
		Apikey string `json:"apikey"`
	} `json:"newyo"`
}

// Load keys from json file
func LoadKeysConfig() Keys_t {
	var k Keys_t

	path, _ := filepath.Abs(keyPath)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&k)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Successful loading of Key Info ........")

	return k
}

var Keys Keys_t

func init() {
	Keys = LoadKeysConfig()
}
