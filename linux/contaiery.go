package linux

type Contaiery struct {
	Configy    Configy
	BandlePath string
}

func NewContaiery(configy Configy, bandlePath string) *Contaiery {
	return &Contaiery{
		Configy:    configy,
		BandlePath: bandlePath,
	}
}
