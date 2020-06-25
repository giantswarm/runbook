package problem

var None = Kind{
	ID: "None",
	Description: "A problem has not been detected.",
}

var Unknown = Kind{
	ID: "Unknown",
	Description: "The problem kind cannot been detected.",
}

type Kind struct {
	ID          string
	Description string
}

func IsFound(problemKinds ...Kind) bool {
	for _, kind := range problemKinds {
		if kind.ID != None.ID || kind.ID != Unknown.ID {
			return true
		}
	}

	return false
}
