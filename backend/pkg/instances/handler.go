package instances

import (
	"incompetent-hosting-provider/backend/pkg/constants"
	"incompetent-hosting-provider/backend/pkg/mq_handler"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type CreatePresetContainerBody struct {
	Preset        int    `json:"preset"`
	ContainerName string `json:"name"`
	Description   string `json:"description"`
}

type ContainerImageDescription struct {
	Tag       string `json:"version"`
	ImageName string `json:"name"`
}

type CreateCustomContainerBody struct {
	Containername string                    `json:"name"`
	Description   string                    `json:"description"`
	Image         ContainerImageDescription `json:"image"`
	EnvVars       map[string]string         `json:"env_vars"`
	Ports         []int                     `json:"ports"`
}

type CreateContainerResponse struct {
	ContainerId string `json:"id"`
}

// godoc
// @Summary 				  	Create container based on preset
//
// @Schemes
// @Description 				Start the container creation flow. This will schedule the creation of said container
// @Tags 						instances
//
// @Security					BearerAuth
//
// @Param request body instances.CreatePresetContainerBody true "query params"
//
// @Success 					202 {object} instances.CreateContainerResponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						503 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /instances/preset [post]
func CreatePresetContainerHandler(c *gin.Context) {
	// Use header set by middleware
	userId := c.Request.Header.Get(constants.USER_ID_HEADER)

	var createRequest CreatePresetContainerBody

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}

	containerId := uuid.NewString()

	err = mq_handler.PublishPresetContainerStartEvent(mq_handler.PresetContainerStartEvent{
		UserId:        userId,
		ContainerUUID: containerId,
		PresetId:      createRequest.Preset,
	})

	if err != nil {
		util.ThrowServiceUnavailableException(c, "Could not schedule container at the current time")
		return
	}

	c.JSON(http.StatusAccepted, CreateContainerResponse{
		ContainerId: containerId,
	})

}

// godoc
// @Summary 				  	Create container based on custom definition
//
// @Schemes
// @Description 				Start the container creation flow. This will schedule the creation of said container
// @Tags 						instances
//
// @Security					BearerAuth
//
// @Param request body instances.CreateCustomContainerBody true "query params"
//
// @Success 					202 {object} instances.CreateContainerResponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						503 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /instances/custom [post]
func CreateCustomContainerHandler(c *gin.Context) {
	// Use header set by middleware
	userId := c.Request.Header.Get(constants.USER_ID_HEADER)

	var createRequest CreateCustomContainerBody

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}

	containerId := uuid.NewString()

	err = mq_handler.PublishCustomContainerStartEvent(mq_handler.CustomContainerStartEvent{
		UserId:            userId,
		ContainerUUID:     containerId,
		ContainerImage:    createRequest.Image.ImageName,
		ContainerImageTag: createRequest.Image.Tag,
		ContainerEnv:      createRequest.EnvVars,
		ContainerPorts:    createRequest.Ports,
	})

	if err != nil {
		util.ThrowServiceUnavailableException(c, "Could not schedule container at the current time")
		return
	}

	c.JSON(http.StatusAccepted, CreateContainerResponse{
		ContainerId: containerId,
	})
}

// godoc
// @Summary 				  	Delete container
//
// @Schemes
// @Description 				Delete container by ID
// @Tags 						instances
//
// @Security					BearerAuth
//
// @Param   containerId     path    string     true        "Container Id"
//
// @Success 					202 {string} string	"accepted"
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						503 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /instances/:containerId [delete]
func DeleteContainerHandler(c *gin.Context) {
	containerId := c.Param("containerId")
	userId := c.Request.Header.Get(constants.USER_ID_HEADER)

	if containerId == "" {
		util.ThrowBadRequestException(c, "No valid containerId passed")
	}

	err := mq_handler.PublishDeleteContainerEvent(mq_handler.DeleteContainerEvent{
		ContainerUUID: containerId,
		UserId:        userId,
	})

	if err != nil {
		util.ThrowServiceUnavailableException(c, "Could not schedule container at the current time")
		return
	}

	c.Status(http.StatusAccepted)
}
