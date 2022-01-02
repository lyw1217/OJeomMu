package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"ojeommu/config"

	"github.com/google/go-querystring/query"
)

/* 카카오 키워드로 장소 검색하기 */
func SearchKeyword(s SearchCond_t) []SearchKeyword_t {

	que := "맛집"
	code := "FD6" // 음식점
	x := s.X
	y := s.Y
	rad, err := strconv.Atoi(s.Radius)
	if err != nil {
		rad = 100
		log.Println("Error radius strconv atoi")
	}

	k := config.Keys.Kakao

	var p = KeywordParam_t{
		Query:             que,
		CategoryGroupCode: code,
		X:                 x,
		Y:                 y,
		Radius:            rad,
		Page:              1,
		Size:              15,
		Sort:              "distance",
	}

	baseUrl, err := url.Parse(KakaoSearchKeywordUrl)
	if err != nil {
		log.Fatal("Malformed URL: ", err.Error())
		return nil
	}

	vals, _ := query.Values(p)
	// Add Query Parameters to the URL
	baseUrl.RawQuery = vals.Encode() // Escape Query Parameters

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "KakaoAK "+k.Rest)

	fmt.Printf("Encoded URL is %q\n", req.URL.String())

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Fatal(str)
	}

	data := make([]SearchKeyword_t, 1, MAX_SEARCH_PAGE)
	err = json.Unmarshal(rspBody, &data[0])
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}

	//fmt.Printf("data : %#v\n", data[0])

	if data[0].Meta.PageableCount > 2 {
		for i := 2; i < data[0].Meta.PageableCount && i < MAX_SEARCH_PAGE; i++ {
			p.Page = i
			baseUrl, err := url.Parse(KakaoSearchKeywordUrl)
			if err != nil {
				log.Fatal("Malformed URL: ", err.Error())
				return nil
			}

			vals, _ := query.Values(p)
			// Add Query Parameters to the URL
			baseUrl.RawQuery = vals.Encode() //params.Encode() // Escape Query Parameters

			req, err := http.NewRequest("GET", baseUrl.String(), nil)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Add("Authorization", "KakaoAK "+k.Rest)

			fmt.Printf("Encoded URL is %q\n", req.URL.String())

			client := &http.Client{}
			rsp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer rsp.Body.Close()

			rspBody, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				str := string(rspBody)
				log.Fatal(str)
			}

			data = append(data, data...)

			err = json.Unmarshal(rspBody, &data[i-1])
			if err != nil {
				log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
			}
		}
	}

	return data
}
