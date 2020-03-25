// Package configs ...
package configs

// Global Config of the application
const (
	Service     = "aws-cleanup"
	Description = "Cleaning AWS"
)

// Dynamic version retrieve with ldflags

// Version represent version of application
var Version string

// GitCommit represent git commit
var GitCommit string

// BuildDate represent date of build
var BuildDate string
