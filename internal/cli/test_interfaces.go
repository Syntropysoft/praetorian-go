package cli

// TestableCommand defines interface for testable commands
type TestableCommand interface {
	Execute(flags interface{}) error
	ValidateFlags(flags interface{}) error
	GetCommandName() string
}

// MockFileWriter for testing file operations
type MockFileWriter struct {
	WriteFileFunc func(filename, content string) error
	Files         map[string]string
}

func (m *MockFileWriter) WriteFile(filename, content string) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, content)
	}
	if m.Files == nil {
		m.Files = make(map[string]string)
	}
	m.Files[filename] = content
	return nil
}

// MockValidator for testing validation operations
type MockValidator struct {
	ValidateFunc func(input interface{}) error
	Validations  []interface{}
}

func (m *MockValidator) Validate(input interface{}) error {
	m.Validations = append(m.Validations, input)
	if m.ValidateFunc != nil {
		return m.ValidateFunc(input)
	}
	return nil
}

// MockConfigGenerator for testing config generation
type MockConfigGenerator struct {
	GenerateFunc func(options interface{}) (string, error)
	Generated    []interface{}
}

func (m *MockConfigGenerator) GenerateConfig(options interface{}) (string, error) {
	m.Generated = append(m.Generated, options)
	if m.GenerateFunc != nil {
		return m.GenerateFunc(options)
	}
	return "# Mock Config\nversion: \"1.0\"", nil
}
