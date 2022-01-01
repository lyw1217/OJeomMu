package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"ojeommu/config"
	"os"

	"github.com/google/go-querystring/query"
)

const kakaoSearchCatUrl string = "https://dapi.kakao.com/v2/local/search/category.json"
const kakaoSearchKeywordUrl string = "https://dapi.kakao.com/v2/local/search/keyword.json"

type CatParam_t struct {
	CategoryGroupCode string `url:"category_group_code"`
	X                 string `url:"x"`
	Y                 string `url:"y"`
	Radius            int    `url:"radius"`
	Page              int    `url:"page"`
	Sort              string `url:"sort"`
}

type CatSameName_t struct {
	Region         []string `json:"region"`
	Keyword        string   `json:"keyword"`
	SelectedRegion string   `json:"selected_region"`
}

type CatMeta_t struct {
	TotalCount    int           `json:"total_count"`
	PageableCount int           `json:"pageable_count"`
	IsEnd         bool          `json:"is_end"`
	SameName      CatSameName_t `json:"same_name"`
}

type CatDocuments_t struct {
	Id                string `json:"id"`
	PlaceName         string `json:"place_name"`
	CategoryName      string `json:"category_name"`
	CategoryGroupCode string `json:"category_group_code"`
	CategoryGroupName string `json:"category_group_name"`
	Phone             string `json:"phone"`
	AddressName       string `json:"address_name"`
	RoadAddressName   string `json:"road_address_name"`
	X                 string `json:"x"`
	Y                 string `json:"y"`
	PlaceUrl          string `json:"place_url"`
	Distance          string `json:"distance"`
}

type SearchCat_t struct {
	Documents []CatDocuments_t `json:"documents"`
	Meta      CatMeta_t        `json:"meta"`
	SameName  CatSameName_t    `json:"same_name"`
}

type KeywordParam_t struct {
	Query             string `url:"query"`
	CategoryGroupCode string `url:"category_group_code"`
	X                 string `url:"x"`
	Y                 string `url:"y"`
	Radius            int    `url:"radius"`
	Page              int    `url:"page"`
	Sort              string `url:"sort"`
}

type KeywordSameName_t struct {
	Region         []string `json:"region"`
	Keyword        string   `json:"keyword"`
	SelectedRegion string   `json:"selected_region"`
}

type KeywordMeta_t struct {
	TotalCount    int           `json:"total_count"`
	PageableCount int           `json:"pageable_count"`
	IsEnd         bool          `json:"is_end"`
	SameName      CatSameName_t `json:"same_name"`
}

type KeywordDocuments_t struct {
	Id                string `json:"id"`
	PlaceName         string `json:"place_name"`
	CategoryName      string `json:"category_name"`
	CategoryGroupCode string `json:"category_group_code"`
	CategoryGroupName string `json:"category_group_name"`
	Phone             string `json:"phone"`
	AddressName       string `json:"address_name"`
	RoadAddressName   string `json:"road_address_name"`
	X                 string `json:"x"`
	Y                 string `json:"y"`
	PlaceUrl          string `json:"place_url"`
	Distance          string `json:"distance"`
}

type SearchKeyword_t struct {
	Documents []KeywordDocuments_t `json:"documents"`
	Meta      KeywordMeta_t        `json:"meta"`
	SameName  KeywordSameName_t    `json:"same_name"`
}

func TestSearchCat(code string, x string, y string, rad int) {

	k := config.Keys.Kakao

	var p = CatParam_t{
		CategoryGroupCode: code,
		X:                 x,
		Y:                 y,
		Radius:            rad,
		Page:              1,
		Sort:              "distance",
	}

	baseUrl, err := url.Parse(kakaoSearchCatUrl)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
	}

	vals, _ := query.Values(p)
	fmt.Println(vals.Encode())
	// Add Query Parameters to the URL
	baseUrl.RawQuery = vals.Encode() //params.Encode() // Escape Query Parameters

	fmt.Printf("Encoded URL is %q\n", baseUrl.String())

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", "KakaoAK "+k.Rest)

	fmt.Println("req.URL.String() : ", req.URL.String())

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

	fmt.Println(string(rspBody))

	//var data []SearchCat_t
	data := make([]SearchCat_t, 1, 10)
	err = json.Unmarshal(rspBody, &data[0])
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}

	if data[0].Meta.PageableCount > 2 {
		for i := 0; i < data[0].Meta.PageableCount; i++ {

		}
	}

	fmt.Printf("data : %#v\n", data[0])

}

func TestSearchKeyword(que string, code string, x string, y string, rad int) {

	k := config.Keys.Kakao

	var p = KeywordParam_t{
		Query:             que,
		CategoryGroupCode: code,
		X:                 x,
		Y:                 y,
		Radius:            rad,
		Page:              1,
		Sort:              "distance",
	}

	baseUrl, err := url.Parse(kakaoSearchKeywordUrl)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
	}

	vals, _ := query.Values(p)
	fmt.Println(vals.Encode())
	// Add Query Parameters to the URL
	baseUrl.RawQuery = vals.Encode() //params.Encode() // Escape Query Parameters

	fmt.Printf("Encoded URL is %q\n", baseUrl.String())

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", "KakaoAK "+k.Rest)

	fmt.Println("req.URL.String() : ", req.URL.String())

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

	fmt.Println(string(rspBody))

	data := make([]SearchKeyword_t, 1, 2)
	err = json.Unmarshal(rspBody, &data[0])
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}

	fmt.Printf("data : %#v\n", data[0])

	if data[0].Meta.PageableCount > 2 {
		for i := 2; i < data[0].Meta.PageableCount && i < 10; i++ {
			p.Page = i
			baseUrl, err := url.Parse(kakaoSearchKeywordUrl)
			if err != nil {
				fmt.Println("Malformed URL: ", err.Error())
				return
			}

			vals, _ := query.Values(p)
			fmt.Println(vals.Encode())
			// Add Query Parameters to the URL
			baseUrl.RawQuery = vals.Encode() //params.Encode() // Escape Query Parameters

			fmt.Printf("Encoded URL is %q\n", baseUrl.String())

			req, err := http.NewRequest("GET", baseUrl.String(), nil)
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}
			req.Header.Add("Authorization", "KakaoAK "+k.Rest)

			fmt.Println("req.URL.String() : ", req.URL.String())

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

			fmt.Println(string(rspBody))

			data = append(data, data...)

			err = json.Unmarshal(rspBody, &data[i-1])
			if err != nil {
				log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
			}
		}
	}
}
