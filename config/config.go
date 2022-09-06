package config

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
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

	k.Naver.Id = os.Getenv("NAVER_ID")
	k.Naver.Secret = os.Getenv("NAVER_SECRET")
	k.Naver.NcpId = os.Getenv("NCP_ID")
	k.Naver.NcpSecret = os.Getenv("NCP_SECRET")

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
