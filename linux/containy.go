package linux

type Containy struct {
	Configy    Configy
	BandlePath string
}

func NewContainy(configy Configy, bandlePath string) *Containy {
	return &Containy{
		Configy:    configy,
		BandlePath: bandlePath,
	}
}
