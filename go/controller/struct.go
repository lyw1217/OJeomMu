package controller

const MAX_SEARCH_PAGE = 10 // 최대 API 호출 횟수
const MAX_SEARCH_DOC = 180 // 최대 매장 개수(45 * 4개 사분면)

/* KAKAO KEYWORD SEARCH */

type KeywordParam_t struct {
	Query             string `url:"query"`
	CategoryGroupCode string `url:"category_group_code"`
	Rect              string `url:"rect"`
	Page              int    `url:"page"`
	Size              int    `url:"size"`
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

/* KAKAO CATEGORY SEARCH */
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

/* SearchHandler */
type SearchCond_t struct {
	Category string `json:"category"`
	X        string `json:"x"`
	Y        string `json:"y"`
	Radius   string `json:"radius"`
}

type Coord_t struct {
	Lat float64
	Lng float64
}

type RectCoord_t struct {
	N Coord_t // north
	S Coord_t // south
	W Coord_t // west
	E Coord_t // east
}

/*
네이버 Search API Local
https://developers.naver.com/docs/serviceapi/search/local/local.md#%EC%A7%80%EC%97%AD
*/
type LocalParam_t struct {
	Query   string `url:"query"`
	Display int    `url:"display"`
	Start   int    `url:"start"` // 검색 시작 위치로 1만 가능
	Sort    string `url:"sort"`  // 정렬 옵션: random(유사도순), comment(카페/블로그 리뷰 개수 순)
}

type LocalDocument_t struct {
	LastBuildDate string `json:"lastBuildDate"`
	Total         int    `json:"total"`
	Start         int    `json:"start"`
	Display       int    `json:"display"`
	Category      string `json:"category"`
	Items         []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Address     string `json:"address"`
		RoadAddress string `json:"roadAddress"`
		Mapx        string `json:"mapx"`
		Mapy        string `json:"mapy"`
	} `json:"items"`
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

/*
네이버 geocode
https://api.ncloud-docs.com/docs/ai-naver-mapsgeocoding-geocode
*/

type GeocodeParam_t struct {
	Query      string `url:"query"`
	Coordinate string `url:"coordinate"` // - 검색 중심 좌표 'lon,lat' 형식으로 입력
	Filter     string `url:"filter"`
	Page       string `url:"page"`
	Count      string `url:"count"`
}

type GeocodeDocument_t struct {
	Status string `json:"status"`
	Meta   struct {
		TotalCount int `json:"totalCount"`
		Page       int `json:"page"`
		Count      int `json:"count"`
	} `json:"meta"`
	Addresses []struct {
		RoadAddress     string `json:"roadAddress"`
		JibunAddress    string `json:"jibunAddress"`
		EnglishAddress  string `json:"englishAddress"`
		AddressElements []struct {
			Types     []string `json:"types"`
			LongName  string   `json:"longName"`
			ShortName string   `json:"shortName"`
			Code      string   `json:"code"`
		} `json:"addressElements"`
		X        string  `json:"x"`
		Y        string  `json:"y"`
		Distance float64 `json:"distance"`
	} `json:"addresses"`
	ErrorMessage string `json:"errorMessage"`
}
