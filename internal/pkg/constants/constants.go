package constants

// Project types
const (
	ProjectTypeAPI     = "api"
	ProjectTypeCLI     = "cli"
	ProjectTypeLibrary = "library"
)

// Database types
const (
	DatabaseMongoDB    = "mongodb"
	DatabasePostgreSQL = "postgresql"
	DatabaseNone       = "none"
)

// Version
const Version = "0.1.0"

// Project type descriptions
var ProjectTypeDescriptions = map[string]string{
	ProjectTypeAPI:     "RESTful API backend with HTTP server",
	ProjectTypeCLI:     "Command-line application",
	ProjectTypeLibrary: "Go library/package",
}

// Database descriptions
var DatabaseDescriptions = map[string]string{
	DatabaseMongoDB:    "MongoDB (NoSQL document database)",
	DatabasePostgreSQL: "PostgreSQL (SQL relational database)",
	DatabaseNone:       "No database integration",
}
