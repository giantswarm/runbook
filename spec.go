package runbook

type Runbook interface {
	// Get the runbook identifier
	GetID() string

	// Get the source URL of the runbook
	GetSourceURL() string

	// Try to apply the runbook
	Apply() error

	// Precondition to check if the runbook can be applied
	Test() (bool, error)
}
