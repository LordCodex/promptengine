package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	apperrors "github.com/LordCodex/promptengine/internal/errors"
	"github.com/spf13/viper"
)

const CurrentVersion = "1.0.0"

type ProjectInfo struct {
	Name             string `mapstructure:"name" json:"name" yaml:"name"`
	Version          string `mapstructure:"version" json:"version" yaml:"version"`
	Description      string `mapstructure:"description" json:"description,omitempty" yaml:"description,omitempty"`
	PromptEnginePath string `mapstructure:"promptengine_path" json:"promptengine_path" yaml:"promptengine_path"`
}

type DocsPaths struct {
	RootDir            string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	AgentsConstitution string `mapstructure:"agents_constitution" json:"agents_constitution" yaml:"agents_constitution"`
	APISpec            string `mapstructure:"api_spec" json:"api_spec" yaml:"api_spec"`
	DatabaseSpec       string `mapstructure:"database_spec" json:"database_spec" yaml:"database_spec"`
	BusinessRules      string `mapstructure:"business_rules" json:"business_rules" yaml:"business_rules"`
}

type CLIConfig struct {
	Verbose bool   `mapstructure:"verbose" json:"verbose" yaml:"verbose"`
	Debug   bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	JSON    bool   `mapstructure:"json" json:"json" yaml:"json"`
	Config  string `mapstructure:"config" json:"config,omitempty" yaml:"config,omitempty"`
}

type AgentProfileConfig struct {
	InstructionFile string `mapstructure:"instruction_file" json:"instruction_file" yaml:"instruction_file"`
	Format          string `mapstructure:"format" json:"format,omitempty" yaml:"format,omitempty"`
}

type AppConfig struct {
	Version string                        `mapstructure:"version" json:"version" yaml:"version"`
	Mode    string                        `mapstructure:"mode" json:"mode" yaml:"mode"`
	Project ProjectInfo                   `mapstructure:"project" json:"project" yaml:"project"`
	Docs    DocsPaths                     `mapstructure:"docs" json:"docs" yaml:"docs"`
	CLI     CLIConfig                     `mapstructure:"cli" json:"cli" yaml:"cli"`
	Agents  map[string]AgentProfileConfig `mapstructure:"agents" json:"agents,omitempty" yaml:"agents,omitempty"`
	Global  map[string]interface{}        `mapstructure:"global" json:"global,omitempty" yaml:"global,omitempty"`
	Plugins map[string]interface{}        `mapstructure:"plugins" json:"plugins,omitempty" yaml:"plugins,omitempty"`
	Org     map[string]interface{}        `mapstructure:"org" json:"org,omitempty" yaml:"org,omitempty"`
	Remote  map[string]interface{}        `mapstructure:"remote" json:"remote,omitempty" yaml:"remote,omitempty"`
}

func DefaultConfig() *AppConfig {
	return &AppConfig{
		Version: CurrentVersion,
		Mode:    "production",
		Project: ProjectInfo{
			Name:             "PromptEngineProject",
			Version:          CurrentVersion,
			PromptEnginePath: "./promptengine",
		},
		Docs: DocsPaths{
			RootDir:            "docs",
			AgentsConstitution: "AGENTS.md",
			APISpec:            "docs/API.md",
			DatabaseSpec:       "docs/Database.md",
			BusinessRules:      "docs/BusinessRules.md",
		},
		Agents:  map[string]AgentProfileConfig{},
		Global:  map[string]interface{}{},
		Plugins: map[string]interface{}{},
		Org:     map[string]interface{}{},
		Remote:  map[string]interface{}{},
	}
}

func (c *AppConfig) Migrate() error {
	if c.Version == "" {
		c.Version = CurrentVersion
	}
	if c.Mode == "" {
		c.Mode = "production"
	}
	if c.Global == nil {
		c.Global = map[string]interface{}{}
	}
	if c.Plugins == nil {
		c.Plugins = map[string]interface{}{}
	}
	if c.Org == nil {
		c.Org = map[string]interface{}{}
	}
	if c.Remote == nil {
		c.Remote = map[string]interface{}{}
	}
	if c.Agents == nil {
		c.Agents = map[string]AgentProfileConfig{}
	}
	return nil
}

type Flags struct {
	Verbose *bool
	Debug   *bool
	JSON    *bool
	Config  string
}

type ConfigLoader struct {
	v *viper.Viper
}

func NewConfigLoader() *ConfigLoader {
	v := viper.New()
	v.SetConfigType("yaml")
	applyDefaults(v)
	v.SetEnvPrefix("PROMPTENGINE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	return &ConfigLoader{v: v}
}

func (l *ConfigLoader) Load(projectPath, globalPath string) (*AppConfig, error) {
	return l.LoadWithFlags(projectPath, globalPath, Flags{})
}

func (l *ConfigLoader) LoadWithFlags(projectPath, globalPath string, flags Flags) (*AppConfig, error) {
	if err := mergeConfigFile(l.v, globalPath); err != nil {
		return nil, configError("load global configuration", err)
	}

	configPath := projectPath
	if flags.Config != "" {
		configPath = flags.Config
	}
	if err := mergeConfigFile(l.v, configPath); err != nil {
		return nil, configError("load project configuration", err)
	}

	if flags.Verbose != nil {
		l.v.Set("cli.verbose", *flags.Verbose)
	}
	if flags.Debug != nil {
		l.v.Set("cli.debug", *flags.Debug)
	}
	if flags.JSON != nil {
		l.v.Set("cli.json", *flags.JSON)
	}
	if flags.Config != "" {
		l.v.Set("cli.config", flags.Config)
	}

	var cfg AppConfig
	if err := l.v.Unmarshal(&cfg); err != nil {
		return nil, configError("decode configuration", err)
	}
	if err := cfg.Migrate(); err != nil {
		return nil, configError("migrate configuration", err)
	}
	return &cfg, nil
}

func configError(action string, err error) error {
	return apperrors.New(apperrors.CategoryConfiguration, apperrors.ExitConfiguration, action, err).
		WithRecommendation("Check the YAML syntax and configuration path, then rerun the command.")
}

func applyDefaults(v *viper.Viper) {
	d := DefaultConfig()
	v.SetDefault("version", d.Version)
	v.SetDefault("mode", d.Mode)
	v.SetDefault("project.name", d.Project.Name)
	v.SetDefault("project.version", d.Project.Version)
	v.SetDefault("project.promptengine_path", d.Project.PromptEnginePath)
	v.SetDefault("docs.root_dir", d.Docs.RootDir)
	v.SetDefault("docs.agents_constitution", d.Docs.AgentsConstitution)
	v.SetDefault("docs.api_spec", d.Docs.APISpec)
	v.SetDefault("docs.database_spec", d.Docs.DatabaseSpec)
	v.SetDefault("docs.business_rules", d.Docs.BusinessRules)
	v.SetDefault("cli.verbose", false)
	v.SetDefault("cli.debug", false)
	v.SetDefault("cli.json", false)
}

func mergeConfigFile(v *viper.Viper, path string) error {
	if path == "" {
		return nil
	}
	expanded, err := expandPath(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(expanded); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	v.SetConfigFile(expanded)
	return v.MergeInConfig()
}

func expandPath(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
	}
	return path, nil
}
