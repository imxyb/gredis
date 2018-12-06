package gredis 

type RangeSpec struct {
	Min float64
	Max float64
	Minex bool 
	Maxex bool 
}