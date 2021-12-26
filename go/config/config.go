package config

import (
	"encoding/json" // https://pkg.go.dev/encoding/json
	"log"
	"os"
	"path/filepath"
)

const (
	keyPath string = "./config/keys.json"
)

type Kakao_t struct {
	Rest  string `json:"rest_api"`
	JS    string `json:"javascript"`
	Admin string `json:"admin"`
}

type Keys_t struct {
	Kakao Kakao_t `json:"kakao"`
}

// Load keys from json file
func LoadKeysConfig() Keys_t {
	var k Keys_t

	path, _ := filepath.Abs(keyPath)
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&k)
	if err != nil {
		log.Println(err)
	}

	log.Print("< SCRAPER > Successful loading of Key Info ........")

	return k
}

var Keys Keys_t

func init() {
	Keys = LoadKeysConfig()
}
