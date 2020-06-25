package problem

var None = Kind{
	ID:          "None",
	Description: "No problem has been identified.",
}

var Unknown = Kind{
	ID:          "Unknown",
	Description: "This kind of problem is unknown. Welcome explorers to a new territory.",
}

type Kind struct {
	ID          string
	Description string
}

func IsFound(kinds ...Kind) bool {
	for _, kind := range kinds {
		if kind.ID != None.ID && kind.ID != Unknown.ID {
			return true
		}
	}

	return false
}
