package subsystems

type ResourceConfig struct {
	CpuLimits    map[string]string
	MemoryLimits map[string]string
}

type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

// 通过不同的subsystem初始化实例，创建资源限制处理数组
var (
	SubsystemsIns = []Subsystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
