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

func IsFound(kind Kind) bool {
	if kind.ID == None.ID || kind.ID == Unknown.ID {
		return false
	} else {
		return true
	}
}
