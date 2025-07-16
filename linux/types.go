package linux

type Configy struct {
	Version       string    `json:"version"`
	RootfsPath    string    `json:"rootfsPath"`
	ContainerName string    `json:"containerName"`
	Cmd           CmdUnit   `json:"cmd"`
	Isolation     Isolation `json:"isolation"`
	UserNS        UserNS    `json:"userns"`
}

type CmdUnit struct {
	Command string `json:"command"`
	CmdArgv string `json:"cmd_argv"`
}

type Isolation struct {
	Cgroups    []CgroupsUint `json:"cgroups"`
	Namespaces []string      `json:"namespaces"`
}

type CgroupsUint struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
	Value    string `json:"value"`
}

type UserNS struct {
	User        UserUnit    `json:"user"`
	UidMappings MappingUint `json:"uidMappings"`
	GidMappings MappingUint `json:"gidMappings"`
}

type UserUnit struct {
	UID int `json:"uid"`
	GID int `json:"gid"`
}

type MappingUint struct {
	ContainerID int `json:"containerID"`
	HostID      int `json:"hostID"`
	Size        int `json:"size"`
}
