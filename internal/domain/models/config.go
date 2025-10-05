package models

import (
	"time"
)

// PraetorianConfig represents the main configuration structure
type PraetorianConfig struct {
	Version      string                    `yaml:"version" json:"version"`
	Files        FilePatterns              `yaml:"files" json:"files"`
	Environments map[string]string         `yaml:"environments" json:"environments"`
	Rules        ValidationRules           `yaml:"rules" json:"rules"`
	Output       OutputConfig              `yaml:"output" json:"output"`
	Performance  PerformanceConfig         `yaml:"performance" json:"performance"`
	Integrations IntegrationConfig         `yaml:"integrations" json:"integrations"`
}

// FilePatterns defines file patterns and exclusions
type FilePatterns struct {
	Include []string `yaml:"include" json:"include"`
	Exclude []string `yaml:"exclude" json:"exclude"`
}

// ValidationRules defines validation rules
type ValidationRules struct {
	Structure  StructureRules  `yaml:"structure" json:"structure"`
	Security   SecurityRules   `yaml:"security" json:"security"`
	Compliance ComplianceRules `yaml:"compliance" json:"compliance"`
}

// StructureRules defines structure validation rules
type StructureRules struct {
	RequiredKeys []string `yaml:"required_keys" json:"required_keys"`
	ForbiddenKeys []string `yaml:"forbidden_keys" json:"forbidden_keys"`
	IgnoreKeys   []string `yaml:"ignore_keys" json:"ignore_keys"`
}

// SecurityRules defines security validation rules
type SecurityRules struct {
	SecretDetection    bool     `yaml:"secret_detection" json:"secret_detection"`
	VulnerabilityScan  bool     `yaml:"vulnerability_scan" json:"vulnerability_scan"`
	PermissionCheck    bool     `yaml:"permission_check" json:"permission_check"`
	CustomPatterns     []string `yaml:"custom_patterns" json:"custom_patterns"`
}

// ComplianceRules defines compliance validation rules
type ComplianceRules struct {
	Standards []string         `yaml:"standards" json:"standards"`
	Policies  []string         `yaml:"policies" json:"policies"`
	Custom    map[string]interface{} `yaml:"custom" json:"custom"`
}

// OutputConfig defines output configuration
type OutputConfig struct {
	Format       string `yaml:"format" json:"format"`
	Colors       bool   `yaml:"colors" json:"colors"`
	Verbose      bool   `yaml:"verbose" json:"verbose"`
	PipelineMode bool   `yaml:"pipeline_mode" json:"pipeline_mode"`
}

// PerformanceConfig defines performance settings
type PerformanceConfig struct {
	Concurrent   bool          `yaml:"concurrent" json:"concurrent"`
	MaxWorkers   int           `yaml:"max_workers" json:"max_workers"`
	Timeout      time.Duration `yaml:"timeout" json:"timeout"`
	MemoryLimit  string        `yaml:"memory_limit" json:"memory_limit"`
}

// IntegrationConfig defines integration settings
type IntegrationConfig struct {
	Notifications NotificationConfig `yaml:"notifications" json:"notifications"`
	Storage       StorageConfig       `yaml:"storage" json:"storage"`
}

// NotificationConfig defines notification settings
type NotificationConfig struct {
	Slack string `yaml:"slack" json:"slack"`
	Teams string `yaml:"teams" json:"teams"`
	Email string `yaml:"email" json:"email"`
}

// StorageConfig defines storage settings
type StorageConfig struct {
	S3  string `yaml:"s3" json:"s3"`
	GCS string `yaml:"gcs" json:"gcs"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *PraetorianConfig {
	return &PraetorianConfig{
		Version: "2.0",
		Files: FilePatterns{
			Include: []string{"configs/*.yaml", "configs/*.json", "configs/*.toml"},
			Exclude: []string{"configs/*.local.*", "configs/*.test.*"},
		},
		Environments: map[string]string{
			"dev":     "configs/dev/*",
			"staging": "configs/staging/*",
			"prod":    "configs/prod/*",
		},
		Rules: ValidationRules{
			Structure: StructureRules{
				RequiredKeys:  []string{"database.host", "api.port"},
				ForbiddenKeys: []string{"debug", "test"},
				IgnoreKeys:    []string{"timestamp", "version"},
			},
			Security: SecurityRules{
				SecretDetection:   true,
				VulnerabilityScan: true,
				PermissionCheck:   true,
			},
			Compliance: ComplianceRules{
				Standards: []string{"PCI_DSS", "GDPR"},
				Policies:  []string{"data_encryption", "access_control"},
			},
		},
		Output: OutputConfig{
			Format:       "text",
			Colors:       true,
			Verbose:      false,
			PipelineMode: false,
		},
		Performance: PerformanceConfig{
			Concurrent:  true,
			MaxWorkers:  4,
			Timeout:     30 * time.Second,
			MemoryLimit: "100MB",
		},
		Integrations: IntegrationConfig{
			Notifications: NotificationConfig{},
			Storage:       StorageConfig{},
		},
	}
}
