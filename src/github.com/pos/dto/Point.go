package dto

import "fmt"

type Point struct {
	Lat float64
	Long float64
}

func NewPoint(lat, long float64) Point {
	point := Point{Lat:lat, Long:long}
	return point
}

func (point Point) toString() string {
	return fmt.Sprintf("%s %s", point.Lat, point.Long)

}