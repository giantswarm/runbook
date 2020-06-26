package sshcommand

type Interface interface {
	// Run the command
	Run(nodeName string, command string) error
}
