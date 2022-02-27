package config

// Config holds the tsfinder configuration
type Config struct {
	BinaryContent string
	Unstealth     bool
	ArgsExec      []string
	SelfRm        bool
	Environ       []string
	Daemon        bool
}
