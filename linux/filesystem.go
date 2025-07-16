package linux

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
)

func (c Contaiery) CreateOverlayFS(containerPath string) func() {
	containerRootFS := fmt.Sprintf("%s/%s", containerPath, "rootfs")

	lowerDir := fmt.Sprintf("%s/%s", c.BandlePath, c.Configy.RootfsPath)
	upperDir := fmt.Sprintf("%s/%s", containerPath, "diff")
	workDir := fmt.Sprintf("%s/%s", containerPath, "work")

	for _, dir := range []string{containerRootFS, upperDir, workDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Panicln(err)
		}
	}

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)
	if err := syscall.Mount("overlay", containerRootFS, "overlay", 0, opts); err != nil {
		log.Panicf("failed to mount overlayfs: %v", err)
	}

	return func() {
		if err := syscall.Unmount(containerRootFS, 0); err != nil {
			log.Printf("failed to unmount container rootfs at %q: %v", containerRootFS, err)
		}

		if err := os.RemoveAll(containerPath); err != nil {
			log.Printf("failed to remove container's /tmp directory: %v", err)
		}
	}
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
