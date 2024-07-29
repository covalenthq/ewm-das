// common/constants.go
package common

// BinaryName is the name of the binary used for both the daemon and CLI tool
const BinaryName = "defaultBinaryName"

// LogFileName generates the log file name based on the binary name
func LogFileName() string {
	return "/var/log/" + BinaryName + ".log"
}
