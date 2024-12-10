package handlers

import (
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"net/http"
)

// GetBrokers godoc
//
//	@Id				GetBrokers
//
//	@Summary		Get all brokers
//	@Description	Gets a list of all brokers.
//	@Tags			Brokers
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		brokers.Broker			"list of brokers"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers [get]
func GetBrokers(w http.ResponseWriter, r *http.Request) {

	result, err := brokers.R().B().GetAll()
	if err != nil {
		render.Error(w, r, err, "Get brokers")
		return
	}

	render.JSON(w, r, result)
}
