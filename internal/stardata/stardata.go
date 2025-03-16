package stardata

type StarData struct {
	Index  int
	Name   string
	SpType string
	Mag    Magnitude

	CI ColorIndex
}

func New(index int, name string, sptype string, mag Magnitude) *StarData {
	ci := NewColorIndex(mag.B-mag.V, mag.U-mag.B)
	return &StarData{index, name, sptype, mag, ci}
}
