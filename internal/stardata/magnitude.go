package stardata

type Magnitude struct {
	U float64
	B float64
	V float64
}

func NewMagnitude(u float64, b float64, v float64) Magnitude {
	return Magnitude{u, b, v}
}
