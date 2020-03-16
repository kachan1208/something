package model

type alias []string
type regions []string
type coordinates []float64
type unlocs []string

type Port struct {
	Name        string
	City        string
	Country     string
	Alias       alias
	Regions     regions
	Coordinates coordinates
	Province    string
	TZ          string
	Unlocs      unlocs
	Code        string
}
