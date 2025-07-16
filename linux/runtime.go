package linux

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
)

func (c Contaiery) ContainerRuntime() {
	log.Printf("[%s] :: container was created\n", c.Configy.ContainerName)
	if err := syscall.Sethostname([]byte(c.Configy.ContainerName)); err != nil {
		log.Fatalf("failed to set a new hostname: %v", err)
	}

	containerRootFS := fmt.Sprintf("/tmp/%s/rootfs", c.Configy.ContainerName)

	if err := syscall.Chroot(containerRootFS); err != nil {
		log.Fatalf("failed to chroot inside a container: %v", err)
	}

	if err := os.Chdir("/"); err != nil {
		log.Fatalf("failed to chdir inside a container: %v", err)
	}

	if err := syscall.Mount("none", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		log.Fatalf("failed to make mount namespace private")
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		log.Fatalf("failed to mount /proc: %v", err)
	}

	if err := os.MkdirAll("/sys/fs/cgroup", 0755); err != nil {
		log.Fatalf("failed to create cgroup dir: %v", err)
	}
	if err := syscall.Mount("none", "/sys/fs/cgroup", "cgroup2", 0, ""); err != nil {
		log.Fatalf("failed to mount cgroup2 /sys/fs/cgroup dir: %v", err)
	}

	cmdArgv := strings.Split(c.Configy.Cmd.CmdArgv, " ")
	log.Printf("[%s] :: running user process (%s %s)\n", c.Configy.ContainerName, c.Configy.Cmd.Command, cmdArgv)
	if err := syscall.Exec(c.Configy.Cmd.Command, cmdArgv, os.Environ()); err != nil {
		log.Fatalf("failed calling user cmd: %v", err)
	}
}
