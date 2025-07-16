package linux

type Configy struct {
	Version    string    `json:"version"`
	Cmd        string    `json:"cmd"`
	RootfsPath string    `json:"rootfsPath"`
	Isolation  Isolation `json:"isolation"`
	UserNS     UserNS    `json:"userns"`
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
