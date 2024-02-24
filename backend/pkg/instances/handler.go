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

type createPresetContainerBody struct {
  Preset string `json:"preset"` 
  ContainerName string `json:"name"`
  Description string `json:"description"`
}

type containerImageDescription struct{
	Tag string `json:"version"`
	ImageName string `json:"name"`
}

type createCustomContainerBody struct {
	Containername string `json:"name"`
	Description string `json:"description"`
	Image containerImageDescription `json:"image"`
	EnvVars map[string]string `json:"env_vars"`
	Ports []string `json:"ports"`
}

type creatContainerResponse  struct{
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
// @Success 					202 {object} user.UserResponse
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

	var createRequest createPresetContainerBody 

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}

	containerId := uuid.NewString()

	err = mq_handler.PublishPresetContainerStartEvent(mq_handler.PresetContainerStartEvent{
		UserId: userId,
		ContainerUUID: containerId,
		PresetId: createRequest.Preset,
	})

	if err != nil{
		util.ThrowServiceUnavailableException(c,"Could not schedule container at the current time")
		return
	}

	c.JSON(http.StatusAccepted, creatContainerResponse{
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
// @Success 					202 {object} user.UserResponse
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

	var createRequest createCustomContainerBody 

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}

	containerId := uuid.NewString()

	err = mq_handler.PublishCustomContainerStartEvent(mq_handler.CustomContainerStartEvent{
		UserId: userId,
		ContainerUUID: containerId,
		ContainerImage: createRequest.Image.ImageName,
		ContainerImageTag: createRequest.Image.Tag,
		ContainerEnv: createRequest.EnvVars,
		ContainerPorts: createRequest.Ports,
	})

	if err != nil{
		util.ThrowServiceUnavailableException(c,"Could not schedule container at the current time")
		return
	}

	c.JSON(http.StatusAccepted, creatContainerResponse{
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
// @Success 					202 {object} user.UserResponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						503 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /instances/:id [delete]

func DeleteContainerHandler(c *gin.Context){
	containerId := c.Param("containerId")
	userId := c.Request.Header.Get(constants.USER_ID_HEADER)


	if containerId == ""{
		util.ThrowBadRequestException(c,"No valid containerId passed")
	}

	err := mq_handler.PublishDeleeteContainerEvent(mq_handler.DeleteContainerEvent{
		ContainerUUID: containerId,	
		UserId: userId,
	})

	if err != nil{
		util.ThrowServiceUnavailableException(c,"Could not schedule container at the current time")
		return
	}

	c.Status(http.StatusAccepted)
}

