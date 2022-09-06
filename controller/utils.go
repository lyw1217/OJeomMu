package controller

import (
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

func ConvRadToDeg(r float64) float64 {

	pi := math.Pi

	return r * (180 / pi)
}

func ConvDegToRad(d float64) float64 {

	pi := math.Pi

	return d * (pi / 180)
}

func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func StrToFloat64(s string) float64 {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("Error, failed StrToFloat64(), s =", s, "err =", err)
		return 0
	}
	return result
}

/* 기준점에서 0, 90, 180, 270도, d 킬로미터 떨어진 곳의 좌표 반환 */
// https://stackoverflow.com/questions/7222382/get-lat-long-given-current-point-distance-and-bearing
func GetRectCoord(dlat float64, dlng float64, d float64) RectCoord_t {

	var result RectCoord_t
	R := 6378.1                                    // Radius of the Earth
	brng := []float64{0, 1.5708, 3.14159, 4.71239} // Bearing is 0, 90, 180, 270 degrees converted to radians.
	lat := ConvDegToRad(dlat)
	lng := ConvDegToRad(dlng)

	for i, b := range brng {
		switch i {
		case 0:
			result.E.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.E.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		case 1:
			result.N.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.N.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		case 2:
			result.W.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.W.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		case 3:
			result.S.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.S.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		}
	}

	return result
}

/* 두 좌표 사이 거리 구하기, 반환값 단위 m(미터) */
// https://stackoverflow.com/questions/18883601/function-to-calculate-distance-between-two-coordinates
func GetDistance(lon1, lat1, lon2, lat2 string) int {
	R := 6378.1

	flon1 := StrToFloat64(lon1)
	flat1 := StrToFloat64(lat2)
	flon2 := StrToFloat64(lon2)
	flat2 := StrToFloat64(lat2)

	dLat := ConvDegToRad(flat2 - flat1)
	dLon := ConvDegToRad(flon2 - flon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(ConvDegToRad(flat1))*math.Cos(ConvDegToRad(flat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c // Distance in km

	return int(d * 1000)
}
