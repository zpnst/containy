package linux

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"
)

var namespacesMap map[string]uintptr = map[string]uintptr{
	"user":    syscall.CLONE_NEWUSER,
	"pid":     syscall.CLONE_NEWPID,
	"network": syscall.CLONE_NEWNET,
	"ipc":     syscall.CLONE_NEWIPC,
	"uts":     syscall.CLONE_NEWUTS,
	"mount":   syscall.CLONE_NEWNS,
	"cgroup":  syscall.CLONE_NEWCGROUP,
}

func (c Containy) GetNamespacesFlags() uintptr {
	var cloneFlags uintptr
	for _, namespace := range c.Configy.Isolation.Namespaces {
		cloneFlags |= namespacesMap[namespace]
	}
	return cloneFlags
}

func (c Containy) GetUIDMappings() []syscall.SysProcIDMap {
	var UIDMappings []syscall.SysProcIDMap
	UIDMappings = append(UIDMappings, syscall.SysProcIDMap{
		ContainerID: c.Configy.UserNS.UidMappings.ContainerID,
		HostID:      c.Configy.UserNS.UidMappings.HostID,
		Size:        c.Configy.UserNS.UidMappings.Size,
	})
	return UIDMappings
}

func (c Containy) GetGIDMappings() []syscall.SysProcIDMap {
	var GIDMappings []syscall.SysProcIDMap
	GIDMappings = append(GIDMappings, syscall.SysProcIDMap{
		ContainerID: c.Configy.UserNS.GidMappings.ContainerID,
		HostID:      c.Configy.UserNS.GidMappings.HostID,
		Size:        c.Configy.UserNS.GidMappings.Size,
	})
	return GIDMappings
}

func (c Containy) ConfigureCgroups() (uintptr, error) {
	baseContainyCgroupPath := "/sys/fs/cgroup/containy.slice"
	cgPath := fmt.Sprintf("%s/%s", baseContainyCgroupPath, c.Configy.ContainerName)
	if err := os.MkdirAll(cgPath, 0755); err != nil {
		return 0, err
	}
	var subtreeControlTypes []string
	for _, cgUnit := range c.Configy.Isolation.Cgroups {
		if !slices.Contains(subtreeControlTypes, fmt.Sprintf("+%s", cgUnit.Type)) {
			subtreeControlTypes = append(subtreeControlTypes, fmt.Sprintf("+%s", cgUnit.Type))
		}
	}

	subtreeControlPath := fmt.Sprintf("%s/cgroup.subtree_control", baseContainyCgroupPath)
	if err := os.WriteFile(subtreeControlPath, []byte(strings.Join(subtreeControlTypes, " ")), 0644); err != nil {
		return 0, err
	}

	for _, cgUnit := range c.Configy.Isolation.Cgroups {
		resourcePath := fmt.Sprintf("%s/%s", cgPath, cgUnit.Resource)
		currentResourceValue := cgUnit.Value
		if err := os.WriteFile(resourcePath, []byte(currentResourceValue), 0644); err != nil {
			return 0, err
		}
	}

	f, err := os.Open(cgPath)
	if err != nil {
		return 0, err
	}
	return f.Fd(), nil
}

func (c Containy) CreateContainer() {
	if os.Geteuid() != 0 {
		log.Fatal("Containy must be run as root user")
	}

	containerPath := fmt.Sprintf("/tmp/%s", c.Configy.ContainerName)
	if err := os.MkdirAll(containerPath, 0755); err != nil {
		log.Fatal(err)
	}

	cgroupFD, err := c.ConfigureCgroups()
	if err != nil {
		log.Panicf("failed to configure cgroups: %v", err)
	}

	umountFunc := c.CreateOverlayFS(containerPath)
	defer umountFunc()

	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Env = append(os.Environ(), "IN_CONTAINER=TRUE")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: c.GetNamespacesFlags(),
		Credential: &syscall.Credential{
			Uid: uint32(c.Configy.UserNS.User.UID),
			Gid: uint32(c.Configy.UserNS.User.GID),
		},
		UidMappings:                c.GetUIDMappings(),
		GidMappings:                c.GetGIDMappings(),
		GidMappingsEnableSetgroups: false,
		UseCgroupFD:                true,
		CgroupFD:                   int(cgroupFD),
	}
	if err := cmd.Run(); err != nil {
		log.Printf("failed to run user cmd: %v\n", err)
	}
}
