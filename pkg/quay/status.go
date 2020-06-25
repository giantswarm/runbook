package quay

import (
	"encoding/json"
	"net/http"

	"github.com/giantswarm/microerror"
)

const (
	QuayURL = "https://8szqd6w4s277.statuspage.io/api/v2/status.json"
)

type QuayStatusResponse struct {
	Status QuayStatus
}

type QuayStatus struct {
	Description string
	Indicator   string
}

func IsQuayDown() (bool, error) {
	response, err := http.Get(QuayURL)
	if err != nil {
		return false, microerror.Mask(err)
	}

	defer response.Body.Close()

	var quayStatusResponse QuayStatusResponse
	err = json.NewDecoder(response.Body).Decode(&quayStatusResponse)
	if err != nil {
		return false, microerror.Mask(err)
	}
	indicator := quayStatusResponse.Status.Indicator
	return indicator == "major" || indicator == "critical", nil
}
