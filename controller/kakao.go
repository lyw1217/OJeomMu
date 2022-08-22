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
// 한 번 조회에 최대 45개 제한, 전체 결과(total_count)만큼의 행을 반환해주지 않음
func GetSearchKeyword(p KeywordParam_t, rad int) ([]KeywordDocuments_t, error) {

	baseUrl, err := url.Parse(KakaoSearchKeywordUrl)
	if err != nil {
		log.Println("Malformed URL: ", err.Error())
		return nil, err
	}

	vals, _ := query.Values(p)
	// Add Query Parameters to the URL
	baseUrl.RawQuery = vals.Encode() // Escape Query Parameters

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	k := config.Keys.Kakao
	req.Header.Add("Authorization", "KakaoAK "+k.Rest)

	//log.Printf("Encoded URL is %q\n", req.URL.String())

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

	data := make([]SearchKeyword_t, 1, MAX_SEARCH_PAGE)
	err = json.Unmarshal(rspBody, &data[0])
	if err != nil {
		log.Printf("Error occured during unmarshaling. Error: %s", err.Error())
		return nil, err
	}

	if data[0].Meta.PageableCount > 2 {
		for i := 2; i <= data[0].Meta.PageableCount && i < MAX_SEARCH_PAGE; i++ {
			p.Page = i

			baseUrl, err := url.Parse(KakaoSearchKeywordUrl)
			if err != nil {
				log.Println("Malformed URL: ", err.Error())
				return nil, err
			}

			vals, _ := query.Values(p)
			// Add Query Parameters to the URL
			baseUrl.RawQuery = vals.Encode() //params.Encode() // Escape Query Parameters

			req, err := http.NewRequest("GET", baseUrl.String(), nil)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			req.Header.Add("Authorization", "KakaoAK "+k.Rest)

			//log.Printf("Encoded URL is %q\n", req.URL.String())

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
				log.Println(err, str)
				return nil, err
			}

			var tmp SearchKeyword_t

			err = json.Unmarshal(rspBody, &tmp)
			if err != nil {
				log.Printf("Error occured during unmarshaling. Error: %s", err.Error())
				return nil, err
			}

			data = append(data, tmp)

			if tmp.Meta.IsEnd {
				break
			}
		}
	}

	fd6_list := make([]KeywordDocuments_t, 0, MAX_SEARCH_DOC)

	for _, d := range data {
		fd6_list = append(fd6_list, d.Documents...)
	}

	//log.Println("len = ", len(fd6_list), " cap = ", cap(fd6_list))

	return fd6_list, nil
}

/* d km 크기의 사분면 4개에서 데이터 조회 */
func RectSearch(cond SearchCond_t) (*KeywordDocuments_t, int, error) {
	// X = longitude, Y = latitude
	// 1. 기준점으로부터 d Km 떨어진 Rect 좌표 가져오기(N,S,W,E)
	lng, err := strconv.ParseFloat(cond.X, 64)
	if err != nil {
		log.Println("Error longitude strconv ParseFloat(lng =", cond.X, ")")
		return nil, 0, err
	}
	lat, err := strconv.ParseFloat(cond.Y, 64)
	if err != nil {
		log.Println("Error latitude strconv ParseFloat(lat =", cond.Y, ")")
		return nil, 0, err
	}

	d, err := strconv.ParseFloat(cond.Radius, 64)
	if err != nil {
		d = 0.3
		log.Println("Error radius strconv ParseFloat(rad =", cond.Radius, "), set default(0.3km)")
	}
	coord := GetRectCoord(lat, lng, d)

	//log.Println(coord.N, coord.S, coord.W, coord.E)

	list := make([][]KeywordDocuments_t, 4)
	for i := range list {
		list[i] = make([]KeywordDocuments_t, 0, MAX_SEARCH_DOC)
	}
	// 2. 가져온 좌표기준 1,2,3,4 사분면 데이터 얻기
	for i := 0; i < 4; i++ {
		var source []string
		switch i {
		case 0:
			source = []string{Float64ToStr(coord.N.Lng), Float64ToStr(coord.N.Lat), Float64ToStr(coord.W.Lng), Float64ToStr(coord.W.Lat)}
		case 1:
			source = []string{Float64ToStr(coord.N.Lng), Float64ToStr(coord.N.Lat), Float64ToStr(coord.E.Lng), Float64ToStr(coord.E.Lat)}
		case 2:
			source = []string{Float64ToStr(coord.S.Lng), Float64ToStr(coord.S.Lat), Float64ToStr(coord.W.Lng), Float64ToStr(coord.W.Lat)}
		case 3:
			source = []string{Float64ToStr(coord.S.Lng), Float64ToStr(coord.S.Lat), Float64ToStr(coord.E.Lng), Float64ToStr(coord.E.Lat)}
		}

		rect := strings.Join(source, ",")
		//log.Println("rect =", rect)

		var p = KeywordParam_t{
			Query:             "맛집",
			CategoryGroupCode: "FD6", // 음식점
			Rect:              rect,
			Page:              1,
			Size:              15,
			Sort:              "accuracy", // 또는 distance
		}
		tmp, err := GetSearchKeyword(p, int(d*1000))
		if err != nil {
			log.Println("Error, Failed GetSearchKeyword()")
			continue
		}
		list[i] = tmp

		//time.Sleep(time.Millisecond * 10)
	}
	result, total_nums := GetCondPlace(list, cond)

	return result, total_nums, nil
}

/* ajax<->gin 카테고리 파싱 */
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

/* 조건에 맞는 음식점 랜덤으로 하나 고르기 */
func GetCondPlace(list [][]KeywordDocuments_t, cond SearchCond_t) (*KeywordDocuments_t, int) {

	category := parseCategory(cond.Category)

	matched_places := make([]KeywordDocuments_t, 0, MAX_SEARCH_DOC)
	if strings.Compare(category, "anything") == 0 {
		for _, l := range list {
			matched_places = append(matched_places, l...)
		}
	} else {
		for _, l := range list {
			for _, place := range l {
				if strings.Contains(place.CategoryName, category) {
					matched_places = append(matched_places, place)
				}
			}
		}
	}

	/* 0 ~ MAX_SEARCH_DOC 사이의 난수 생성 */
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(MAX_SEARCH_DOC)
	if len(matched_places) == 0 {
		log.Println("Error, len(matched_place) is zero(0)")
		return nil, 0
	}

	//log.Println("rand num =", n, "len(matched_places) =", len(matched_places))

	result := &matched_places[n%len(matched_places)]

	return result, len(matched_places)
}
