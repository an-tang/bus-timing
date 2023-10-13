package location

import (
	"math"
)

type Location struct {
	Lat float64
	Lng float64
}

func CalculateStraightLine(location1, location2 Location) float64 {
	earthRadius := 6371.0 // km (change this constant to get miles)
	dLat := (location2.Lat - location1.Lat) * (math.Pi / 180)
	dLng := (location2.Lng - location1.Lng) * (math.Pi / 180)

	haversine := (math.Sin(dLat/2) * math.Sin(dLat/2)) + (math.Cos(location1.Lat*math.Pi/180) * math.Cos(location2.Lat*math.Pi/180) * math.Sin(dLng/2) * math.Sin(dLng/2))
	distance := 2 * earthRadius * math.Atan2(math.Sqrt(haversine), math.Sqrt(1-haversine))
	return math.Round(distance * 1000)
}

func IsPointBetween(A, B, X Location) bool {
	// check straightLine A, B, X
	AB := CalculateStraightLine(A, B)
	AX := CalculateStraightLine(A, X)
	BX := CalculateStraightLine(B, X)
	switch {
	case AX+BX == AB:
		return true
	case math.Abs(AX-BX) == AB:
		return false
	}

	// incase A, B, X are not in straight line
	// vector AX = (X(lat) - A(lat) , X(lng) - A(lng))
	// vector XB = (B(lat) - X(lat) , B(lng) - X(lng))
	// if vectorAX * vector XB < 0 => X between A & B
	// if vectorAX * vector XB > 0 => X is not between A & B
	productScalar := (X.Lat-A.Lat)*(B.Lat-X.Lat) + (X.Lng-A.Lng)*(B.Lng-X.Lng)
	switch {
	case productScalar < 0:
		return true
	case productScalar > 0:
		return false
	}

	return false
}
