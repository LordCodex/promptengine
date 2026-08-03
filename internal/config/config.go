package config

import (
	"strings"

	"github.com/spf13/viper"
)

// ProjectInfo schema parameters
type ProjectInfo struct {
	Name             string `mapstructure:"name"`
	Version          string `mapstructure:"version"`
	Description      string `mapstructure:"description"`
	PromptEnginePath string `mapstructure:"promptengine_path"`
}

// DocsPaths schema parameters
type DocsPaths struct {
	RootDir            string `mapstructure:"root_dir"`
	AgentsConstitution string `mapstructure:"agents_constitution"`
	APISpec            string `mapstructure:"api_spec"`
	DatabaseSpec       string `mapstructure:"database_spec"`
	BusinessRules      string `mapstructure:"business_rules"`
}

// AppConfig represents consolidated configurations
type AppConfig struct {
	Project ProjectInfo `mapstructure:"project"`
	Docs    DocsPaths   `mapstructure:"docs"`
}

// ConfigLoader wires Viper parameter precedence
type ConfigLoader struct {
	viperInst *viper.Viper
}

func NewConfigLoader() *ConfigLoader {
	v := viper.New()
	
	// Set defaults
	v.SetDefault("project.name", "PromptEngineProject")
	v.SetDefault("project.version", "1.0.0")
	v.SetDefault("project.promptengine_path", "./promptengine")
	v.SetDefault("docs.root_dir", "docs")
	v.SetDefault("docs.agents_constitution", "AGENTS.md")
	v.SetDefault("docs.api_spec", "docs/API.md")
	v.SetDefault("docs.database_spec", "docs/Database.md")
	v.SetDefault("docs.business_rules", "docs/BusinessRules.md")

	// Set environment configurations bindings
	v.SetEnvPrefix("PROMPTENGINE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return &ConfigLoader{viperInst: v}
}

func (l *ConfigLoader) Load(localPath, globalPath string) (*AppConfig, error) {
	// 1. Load global config file if exists
	if globalPath != "" {
		l.viperInst.SetConfigFile(globalPath)
		_ = l.viperInst.MergeInConfig() // ignore missing global config files
	}

	// 2. Load local config file if exists
	if localPath != "" {
		l.viperInst.SetConfigFile(localPath)
		_ = l.viperInst.MergeInConfig() // ignore missing local config files
	}

	var cfg AppConfig
	if err := l.viperInst.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
