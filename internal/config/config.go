package config

// ProjectConfig represents the project-level configuration
type ProjectConfig struct {
	Platform     string `yaml:"platform"`
	Architecture string `yaml:"architecture"`
	PackageName  string `yaml:"packageName"` // For Android/Java-based projects
}

// ModuleConfig represents the module definition
type ModuleConfig struct {
	Name      string                    `yaml:"name"`
	Type      string                    `yaml:"type"`
	Structure map[string]PlatformConfig `yaml:"structure"`
}

// PlatformConfig holds the file structure for a specific platform
type PlatformConfig struct {
	Data         []string `yaml:"data"`
	Domain       []string `yaml:"domain"`
	Presentation []string `yaml:"presentation"`
}

// Config holds the full configuration
type Config struct {
	Project      ProjectConfig       `yaml:"project"`
	Module       ModuleConfig        `yaml:"module"`
	Dependencies map[string][]string `yaml:"dependencies"`
	Options      Options             `yaml:"options"`
}

type Options struct {
	IncludeTests      bool   `yaml:"includeTests"`
	IncludeDI         bool   `yaml:"includeDI"`
	IncludeNetworking bool   `yaml:"includeNetworking"`
	IncludeLocalDB    bool   `yaml:"includeLocalDB"`
	StateManagement   string `yaml:"stateManagement"`
}
