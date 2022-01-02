package controller

const MAX_SEARCH_PAGE = 10

/* KAKAO KEYWORD SEARCH */

type KeywordParam_t struct {
	Query             string `url:"query"`
	CategoryGroupCode string `url:"category_group_code"`
	X                 string `url:"x"`
	Y                 string `url:"y"`
	Radius            int    `url:"radius"`
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
type SearchCond struct {
	Query  string `json:"query"`
	Code   string `json:"code"`
	X      string `json:"x"`
	Y      string `json:"y"`
	Radius string `json:"radius"`
}
