package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"ojeommu/config"

	"github.com/google/go-querystring/query"
)

/* 네이버 검색 지역, https://developers.naver.com/docs/serviceapi/search/local/local.md#%EC%A7%80%EC%97%AD*/
func GetNaverLocal(p LocalParam_t) (*LocalDocument_t, error) {

	append_qry := []string{"맛집", "한식", "일식", "중식", "분식", "아시아음식", "도시락", "육류", "치킨", "패스트푸드", "술집"}
	sort_qry := []string{"random", "comment"}

	data := []LocalDocument_t{}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(append_qry))

	for j := 0; j < len(sort_qry); j++ {
		baseUrl, err := url.Parse(NaverSearchLocalUrl)
		if err != nil {
			log.Println("Malformed URL: ", err.Error())
			return nil, err
		}

		p.Query = fmt.Sprintf("%s %s", p.Query, append_qry[n])
		p.Sort = sort_qry[j]
		vals, _ := query.Values(p)
		// Add Query Parameters to the URL
		baseUrl.RawQuery = vals.Encode() // Escape Query Parameters

		req, err := http.NewRequest("GET", baseUrl.String(), nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		k := config.Keys.Naver
		req.Header.Add("X-Naver-Client-Id", k.Id)
		req.Header.Add("X-Naver-Client-Secret", k.Secret)

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rsp.Body.Close()

		rspBody, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			str := string(rspBody)
			log.Println(str)
			return nil, err
		}

		var tmp LocalDocument_t
		err = json.Unmarshal(rspBody, &tmp)
		if err != nil {
			log.Printf("Error occured during unmarshaling. Error: %s", err.Error())
			return nil, err
		}

		if tmp.ErrorCode == "" {
			data = append(data, tmp)
		}

		time.Sleep(time.Millisecond * 2)
	}

	n = rand.Intn(len(sort_qry))

	return &data[n], nil
}
