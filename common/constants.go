// common/constants.go
package common

// BinaryName is the name of the binary used for both the daemon and CLI tool
const BinaryName = "defaultBinaryName"

// LogFileName generates the log file name based on the binary name
func LogFileName() string {
	return "/var/log/" + BinaryName + ".log"
}

// Version is the version of the binary
const Version = "v0.1.0"

// GitCommit is the commit hash of the binary
const GitCommit = "00000000"
