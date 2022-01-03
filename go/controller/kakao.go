package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ojeommu/config"

	"github.com/google/go-querystring/query"
)

/* 카카오 키워드로 장소 검색하기 */
func SearchKeyword(s SearchCond_t) ([]KeywordDocuments_t, *SearchCond_t, error) {

	rad, err := strconv.Atoi(s.Radius)
	if err != nil {
		rad = 100
		log.Println("Error radius strconv atoi, set default(100m)")
	}

	var p = KeywordParam_t{
		Query:             "맛집",
		CategoryGroupCode: "FD6", // 음식점
		X:                 s.X,
		Y:                 s.Y,
		Radius:            rad,
		Page:              1,
		Size:              15,
		Sort:              "accuracy", // 또는 distance
	}

	baseUrl, err := url.Parse(KakaoSearchKeywordUrl)
	if err != nil {
		log.Println("Malformed URL: ", err.Error())
		return nil, nil, err
	}

	vals, _ := query.Values(p)
	// Add Query Parameters to the URL
	baseUrl.RawQuery = vals.Encode() // Escape Query Parameters

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	k := config.Keys.Kakao
	req.Header.Add("Authorization", "KakaoAK "+k.Rest)

	//log.Printf("Encoded URL is %q\n", req.URL.String())

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Println(str)
		return nil, nil, err
	}

	data := make([]SearchKeyword_t, 1, MAX_SEARCH_PAGE)
	err = json.Unmarshal(rspBody, &data[0])
	if err != nil {
		log.Printf("Error occured during unmarshaling. Error: %s", err.Error())
		return nil, nil, err
	}

	if data[0].Meta.PageableCount > 2 {
		for i := 2; i <= data[0].Meta.PageableCount && i < MAX_SEARCH_PAGE; i++ {
			p.Page = i

			baseUrl, err := url.Parse(KakaoSearchKeywordUrl)
			if err != nil {
				log.Println("Malformed URL: ", err.Error())
				return nil, nil, err
			}

			vals, _ := query.Values(p)
			// Add Query Parameters to the URL
			baseUrl.RawQuery = vals.Encode() //params.Encode() // Escape Query Parameters

			req, err := http.NewRequest("GET", baseUrl.String(), nil)
			if err != nil {
				log.Println(err)
				return nil, nil, err
			}
			req.Header.Add("Authorization", "KakaoAK "+k.Rest)

			//log.Printf("Encoded URL is %q\n", req.URL.String())

			client := &http.Client{}
			rsp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return nil, nil, err
			}
			defer rsp.Body.Close()

			rspBody, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				str := string(rspBody)
				log.Println(err, str)
				return nil, nil, err
			}

			var tmp SearchKeyword_t

			err = json.Unmarshal(rspBody, &tmp)
			if err != nil {
				log.Printf("Error occured during unmarshaling. Error: %s", err.Error())
				return nil, nil, err
			}

			data = append(data, tmp)

			if tmp.Meta.IsEnd {
				break
			}

			time.Sleep(time.Millisecond * 10)
		}
	}

	//log.Println("len = ", len(data), " cap = ", cap(data))

	/* logging pretty json
	for i := 0; i < len(data); i++ {
		body, _ := json.Marshal(data[i].Documents)
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "\t")
		if err != nil {
			log.Println("JSON parse error: ", err)
			return nil
		}

		log.Println("CSP Violation:", prettyJSON.String())
	}
	*/

	fd6_list := make([]KeywordDocuments_t, 0, MAX_SEARCH_DOC)

	for _, d := range data {
		fd6_list = append(fd6_list, d.Documents...)
	}

	return fd6_list, &s, nil
}

func parseCategory(c string) string {

	switch c {
	case "anything":
		return "anything"
	case "korea":
		return "한식"
	case "china":
		return "중식"
	case "japan":
		return "일식"
	case "western":
		return "양식"
	case "flour":
		return "분식"
	case "asia":
		return "아시아음식"
	case "lunchbox":
		return "도시락"
	case "meat":
		return "육류,고기"
	case "chicken":
		return "치킨"
	case "fastfood":
		return "패스트푸드"
	case "bar":
		return "술집"
	default:
		return "anything"
	}
}

/* 조건에 맞는 음식점 고르기 */
func GetCondPlace(list []KeywordDocuments_t, cond *SearchCond_t) *KeywordDocuments_t {

	category := parseCategory(cond.Category)
	//log.Println("category =", category)

	matched_place := make([]KeywordDocuments_t, 0, len(list))
	if strings.Compare(category, "anything") == 0 {
		matched_place = list
	} else {
		for _, place := range list {
			if strings.Contains(place.CategoryName, category) {
				matched_place = append(matched_place, place)
			}
		}
	}

	/* 1~100 사이의 난수 생성 */
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	if len(matched_place) == 0 {
		log.Println("Error, len(matched_place) is zero(0)")
		return nil
	}
	n = n % len(matched_place)

	return &matched_place[n]
}
