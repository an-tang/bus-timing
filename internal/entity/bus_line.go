package entity

type BusLine struct {
	ID           string
	FullName     string
	ShortName    string
	Origin       string
	BusLinePaths []BusLinePath
}
