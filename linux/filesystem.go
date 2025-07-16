package linux

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func CreateOverlayFS() {

}

func ParseConfigy(configyPath string) (*Configy, error) {
	f, err := os.Open(configyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open configy file, %v", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read configy file, %v", err)
	}

	var configy Configy
	if err := json.Unmarshal(b, &configy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configy file, %v", err)
	}
	return &configy, nil
}
