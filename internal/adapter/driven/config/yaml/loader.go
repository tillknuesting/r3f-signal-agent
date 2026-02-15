package yaml

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"r3f-trends/internal/domain/entity"
)

type Config struct {
	ActiveProfile string          `yaml:"active_profile"`
	Server        ServerConfig    `yaml:"server"`
	Scheduler     SchedulerConfig `yaml:"scheduler"`
	Chrome        ChromeConfig    `yaml:"chrome"`
	LLM           LLMConfig       `yaml:"llm"`
	Storage       StorageConfig   `yaml:"storage"`
	Logging       LoggingConfig   `yaml:"logging"`
}

type ServerConfig struct {
	Host string     `yaml:"host"`
	Port int        `yaml:"port"`
	Auth AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	Enabled bool   `yaml:"enabled"`
	Type    string `yaml:"type"`
	APIKey  string `yaml:"api_key"`
}

type SchedulerConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval string `yaml:"interval"`
	Timezone string `yaml:"timezone"`
}

type ChromeConfig struct {
	Headless       bool     `yaml:"headless"`
	Timeout        string   `yaml:"timeout"`
	BlockResources []string `yaml:"block_resources"`
}

type LLMConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	APIKey   string `yaml:"api_key"`
	BaseURL  string `yaml:"base_url"`
}

type StorageConfig struct {
	Type     string `yaml:"type"`
	BasePath string `yaml:"base_path"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type ConfigLoader struct {
	configPath string
}

func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{configPath: configPath}
}

func (l *ConfigLoader) Load() (*Config, error) {
	data, err := os.ReadFile(filepath.Join(l.configPath, "config.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	l.expandEnv(&cfg)

	return &cfg, nil
}

func (l *ConfigLoader) expandEnv(cfg *Config) {
	cfg.LLM.APIKey = os.ExpandEnv(cfg.LLM.APIKey)
}

type SourceLoader struct {
	sourcesPath string
}

func NewSourceLoader(sourcesPath string) *SourceLoader {
	return &SourceLoader{sourcesPath: sourcesPath}
}

func (l *SourceLoader) LoadAll(ctx context.Context) ([]*entity.Source, error) {
	return l.LoadByProfile(ctx, "tech")
}

func (l *SourceLoader) LoadByProfile(ctx context.Context, profile string) ([]*entity.Source, error) {
	profilePath := filepath.Join(l.sourcesPath, profile)

	files, err := filepath.Glob(filepath.Join(profilePath, "*.yaml"))
	if err != nil {
		return nil, err
	}

	var sources []*entity.Source

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var fileStruct struct {
			Sources []entity.SourceDTO `yaml:"sources"`
		}

		if err := yaml.Unmarshal(data, &fileStruct); err != nil {
			continue
		}

		for _, dto := range fileStruct.Sources {
			sources = append(sources, entity.SourceFromDTO(&dto))
		}
	}

	return sources, nil
}

type ProfileLoader struct {
	profilesPath string
}

func NewProfileLoader(profilesPath string) *ProfileLoader {
	return &ProfileLoader{profilesPath: profilesPath}
}

func (l *ProfileLoader) Load(ctx context.Context, name string) (*entity.Profile, error) {
	data, err := os.ReadFile(filepath.Join(l.profilesPath, name+".yaml"))
	if err != nil {
		return nil, err
	}

	var dto entity.ProfileDTO
	if err := yaml.Unmarshal(data, &dto); err != nil {
		return nil, err
	}

	return entity.ProfileFromDTO(&dto), nil
}

func (l *ProfileLoader) LoadAll(ctx context.Context) ([]*entity.Profile, error) {
	files, err := filepath.Glob(filepath.Join(l.profilesPath, "*.yaml"))
	if err != nil {
		return nil, err
	}

	var profiles []*entity.Profile

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var dto entity.ProfileDTO
		if err := yaml.Unmarshal(data, &dto); err != nil {
			continue
		}

		profiles = append(profiles, entity.ProfileFromDTO(&dto))
	}

	return profiles, nil
}
