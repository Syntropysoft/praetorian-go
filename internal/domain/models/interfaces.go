package models

import (
	"time"
)

// ConfigParser defines the interface for parsing configuration files
type ConfigParser interface {
	CanHandle(filename string) bool
	Parse(content []byte) (*ConfigData, error)
	GetSupportedExtensions() []string
}

// ValidationRule defines the interface for validation rules
type ValidationRule interface {
	ID() string
	Name() string
	Description() string
	Validate(data *ConfigData) ValidationResult
	Severity() SeverityLevel
}

// AuditEngine defines the interface for audit engines
type AuditEngine interface {
	Type() AuditType
	Run(ctx AuditContext) AuditResult
	GetMetrics() AuditMetrics
}

// OutputFormatter defines the interface for output formatters
type OutputFormatter interface {
	Format(result ValidationResult) ([]byte, error)
	GetContentType() string
	SupportsFormat(format string) bool
}

// ConfigData represents parsed configuration data
type ConfigData struct {
	Filename  string                 `json:"filename"`
	Format    string                 `json:"format"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// ValidationResult represents the result of a validation
type ValidationResult struct {
	Success      bool                    `json:"success"`
	Errors       []ValidationError       `json:"errors,omitempty"`
	Warnings     []ValidationWarning     `json:"warnings,omitempty"`
	Summary      ValidationSummary       `json:"summary"`
	Metadata     map[string]interface{}  `json:"metadata,omitempty"`
	Timestamp    time.Time               `json:"timestamp"`
	Duration     time.Duration           `json:"duration"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Severity    SeverityLevel `json:"severity"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
	Column      int    `json:"column,omitempty"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Severity    SeverityLevel `json:"severity"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
	Column      int    `json:"column,omitempty"`
}

// ValidationSummary represents validation summary statistics
type ValidationSummary struct {
	TotalKeys        int `json:"total_keys"`
	MissingKeys      int `json:"missing_keys"`
	ExtraKeys        int `json:"extra_keys"`
	ValueDifferences int `json:"value_differences"`
	SecurityIssues   int `json:"security_issues"`
	ComplianceIssues int `json:"compliance_issues"`
}

// AuditResult represents the result of an audit
type AuditResult struct {
	Type        AuditType              `json:"type"`
	Success     bool                   `json:"success"`
	Score       float64                `json:"score"`
	Grade       string                 `json:"grade"`
	Issues      []AuditIssue           `json:"issues,omitempty"`
	Metrics     AuditMetrics           `json:"metrics"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Duration    time.Duration          `json:"duration"`
}

// AuditIssue represents an audit issue
type AuditIssue struct {
	Type        string        `json:"type"`
	Severity    SeverityLevel `json:"severity"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	File        string        `json:"file,omitempty"`
	Line        int           `json:"line,omitempty"`
	Column      int           `json:"column,omitempty"`
	Recommendation string     `json:"recommendation,omitempty"`
}

// AuditMetrics represents audit metrics
type AuditMetrics struct {
	TotalChecks     int     `json:"total_checks"`
	PassedChecks    int     `json:"passed_checks"`
	FailedChecks    int     `json:"failed_checks"`
	WarningChecks   int     `json:"warning_checks"`
	CriticalIssues  int     `json:"critical_issues"`
	SecurityIssues  int     `json:"security_issues"`
	ComplianceIssues int    `json:"compliance_issues"`
	PerformanceScore float64 `json:"performance_score"`
}

// AuditContext represents the context for an audit
type AuditContext struct {
	Files       []string               `json:"files"`
	Environment string                 `json:"environment"`
	Config      map[string]interface{} `json:"config"`
	Options     map[string]interface{} `json:"options"`
}

// SeverityLevel represents the severity of an issue
type SeverityLevel string

const (
	SeverityInfo     SeverityLevel = "info"
	SeverityLow      SeverityLevel = "low"
	SeverityMedium   SeverityLevel = "medium"
	SeverityHigh     SeverityLevel = "high"
	SeverityCritical SeverityLevel = "critical"
)

// AuditType represents the type of audit
type AuditType string

const (
	AuditTypeSecurity    AuditType = "security"
	AuditTypeCompliance  AuditType = "compliance"
	AuditTypePerformance AuditType = "performance"
	AuditTypeAll         AuditType = "all"
)
