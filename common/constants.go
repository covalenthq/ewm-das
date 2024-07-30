package common

// BinaryName is the name of the binary used for both the daemon and CLI tool
var BinaryName = "defaultBinaryName"

// Version is the version of the binary
var Version = "v0.1.0"

// GitCommit is the commit hash of the binary
var GitCommit = "00000000"

// LogFileName generates the log file name based on the binary name
func LogFileName() string {
	return "/var/log/" + BinaryName + ".log"
}
