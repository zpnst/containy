package linux

type Containy struct {
	Configy    Configy
	BundlePath string
}

func NewContainy(configy Configy, bundlePath string) *Containy {
	return &Containy{
		Configy:    configy,
		BundlePath: bundlePath,
	}
}
