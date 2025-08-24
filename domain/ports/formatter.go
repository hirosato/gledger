package ports

// Formatter defines the interface for formatting output
type Formatter interface {
	// Format converts data to the appropriate output format
	Format(data interface{}) (string, error)
}