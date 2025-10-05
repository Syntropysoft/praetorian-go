package cli

// CommandExecutor defines the interface for command execution
type CommandExecutor interface {
	Execute(args []string) error
}

// FlagExtractor defines the interface for extracting flags from commands
type FlagExtractor interface {
	ExtractFlags() (interface{}, error)
}

// Validator defines the interface for validation operations
type Validator interface {
	Validate(input interface{}) error
}

// ConfigGenerator defines the interface for configuration generation
type ConfigGenerator interface {
	GenerateConfig(options interface{}) (string, error)
}

// FileWriter defines the interface for file writing operations
type FileWriter interface {
	WriteFile(filename, content string) error
}

// OutputFormatter defines the interface for formatting output
type OutputFormatter interface {
	Format(data interface{}) (string, error)
}

// CommandHandler defines the interface for command handling
type CommandHandler interface {
	Handle(command string, args []string) error
}

// ValidationService defines the interface for validation services
type ValidationService interface {
	ValidateConfig(configPath string) error
	ValidateFiles(files []string) error
}

// AuditService defines the interface for audit services
type AuditService interface {
	RunAudit(auditType string, configPath string) error
	RunSecurityAudit(configPath string) error
	RunComplianceAudit(configPath string) error
	RunPerformanceAudit(configPath string) error
}

// ConfigService defines the interface for configuration services
type ConfigService interface {
	CreateConfig(options interface{}) error
	LoadConfig(configPath string) (interface{}, error)
	SaveConfig(config interface{}, configPath string) error
}
