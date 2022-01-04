package controller

import (
	"math"
)

func ConvRadToDeg(r float64) float64 {

	pi := math.Pi

	return r * (180 / pi)
}

func ConvDegToRad(d float64) float64 {

	pi := math.Pi

	return d * (pi / 180)
}

type Coord_t struct {
	Lat float64
	Lng float64
}

type RectCoord_t struct {
	Up    Coord_t
	Down  Coord_t
	Left  Coord_t
	Right Coord_t
}

// 기준점에서 0, 90, 180, 270도, d 킬로미터 떨어진 곳의 좌표 반환
func GetRectCoord(dlat float64, dlng float64, d float64) RectCoord_t {

	var result RectCoord_t
	R := 6378.1                                    // Radius of the Earth
	brng := []float64{0, 1.5708, 3.14159, 4.71239} // Bearing is 0, 90, 180, 270 degrees converted to radians.
	lat := ConvDegToRad(dlat)
	lng := ConvDegToRad(dlng)

	for i, b := range brng {
		switch i {
		case 0:
			result.Right.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.Right.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		case 1:
			result.Up.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.Up.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		case 2:
			result.Left.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.Left.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		case 3:
			result.Down.Lat = ConvRadToDeg(
				math.Asin(math.Sin(lat)*math.Cos(d/R) +
					math.Cos(lat)*math.Sin(d/R)*math.Cos(b)))
			result.Down.Lng = ConvRadToDeg(
				lng + math.Atan2(math.Sin(b)*math.Sin(d/R)*math.Cos(lat),
					math.Cos(d/R)-math.Sin(lat)*math.Sin(lat)))
		}
	}

	return result
}
