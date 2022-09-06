package config

import (
	"encoding/json" // https://pkg.go.dev/encoding/json
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	keyPath   string = "./config/keys.json"
	ServerCrt string = "/etc/letsencrypt/live/mumeog.site/fullchain.pem"
	ServerKey string = "/etc/letsencrypt/live/mumeog.site/privkey.pem"
)

const (
	loggingPath string = "./config/logging.json"
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

/*
log.Trace("Something very low level.")
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// Calls os.Exit(1) after logging
log.Fatal("Bye.")
// Calls panic() after logging
log.Panic("I'm bailing.")
*/
// setup logger
func SetupLogger() {
	path, _ := filepath.Abs(loggingPath)
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var l *lumberjack.Logger

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&l)
	if err != nil {
		log.Println(err)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err, "Err. Failed to get Hostname")
	}
	l.Filename = fmt.Sprintf(l.Filename, hostname)

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, l) // Stderr와 파일에 동시  출력
	gin.DefaultWriter = io.MultiWriter(os.Stdout, l)

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z // or RFC3339
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
	log.SetReportCaller(true) // 해당 이벤트 발생 시 함수, 파일명 표기
	log.Error("Successful Logger setup ...............")
}

var Keys Keys_t

func init() {
	Keys = LoadKeysConfig()
	SetupLogger()
}
