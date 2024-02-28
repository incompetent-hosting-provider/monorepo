package service

import (
	db_presets "incompetent-hosting-provider/backend/pkg/db/tables/presets"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type preset struct {
	Name        string `json:"name"`
	PresetId    int    `json:"id"`
	Description string `json:"description"`
}

type PresetListResponse struct {
	Presets []preset `json:"presets"`
}

// godoc
// @Summary 					Get list of presets
//
// @Schemes
// @Description 				Get full list of available presets
// @Tags 						service
//
// @Security					BearerAuth
//
// @Success 					200 {object} service.PresetListResponse
//
// @Failure						401	{object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /service/available-presets [get]
func GetAvailablePresetsHandler(c *gin.Context) {

	res, err := db_presets.GetAllPresets()

	if err != nil {
		util.ThrowInternalServerErrorException(c, "Could not get list")
		return
	}

	c.JSON(http.StatusOK, serializePresetListResponse(res))
}
