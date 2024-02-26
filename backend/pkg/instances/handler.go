package instances

import (
	"incompetent-hosting-provider/backend/pkg/constants"
	db_instances "incompetent-hosting-provider/backend/pkg/db/tables/instances"
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
	ContainerName string                    `json:"name"`
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
	userSub := c.Request.Header.Get(constants.USER_ID_HEADER)

	var createRequest CreatePresetContainerBody

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}

	containerId := uuid.NewString()

	err = mq_handler.PublishPresetContainerStartEvent(mq_handler.PresetContainerStartEvent{
		UserId:        userSub,
		ContainerUUID: containerId,
		PresetId:      createRequest.Preset,
	})

	if err != nil {
		util.ThrowServiceUnavailableException(c, "Could not schedule container at the current time")
		return
	}

	db_instances.InsertInstance(db_instances.InstancesTable{
		UserSub: userSub,
		ContainerUUID: containerId,
		ContainerPorts: []int{},
		ContainerDescription: createRequest.Description,
		ContainerName: createRequest.ContainerName,
		ImageName: "asda",
		ImageTag: "ajkdhas",
		InstanceStatus: db_instances.STATUS_VALUE_SCHEDULED,
	})

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
	userSub := c.Request.Header.Get(constants.USER_ID_HEADER)

	var createRequest CreateCustomContainerBody

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}

	containerId := uuid.NewString()

	err = mq_handler.PublishCustomContainerStartEvent(mq_handler.CustomContainerStartEvent{
		UserId:            userSub,
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


	db_instances.InsertInstance(db_instances.InstancesTable{
		UserSub: userSub,
		ContainerUUID: containerId,
		ContainerPorts: createRequest.Ports,
		ContainerDescription: createRequest.Description,
		ContainerName: createRequest.ContainerName,
		ImageName: createRequest.Image.ImageName,
		ImageTag: createRequest.Image.Tag,
		InstanceStatus: db_instances.STATUS_VALUE_SCHEDULED,
	})

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
		util.ThrowServiceUnavailableException(c, "Could  schedule container deletion at the current time")
		return
	}

	err = db_instances.DeleteInstanceById(userId, containerId)
	if err != nil{
		util.ThrowInternalServerErrorException(c, "Could not delete entry at this time")
	}

	c.Status(http.StatusAccepted)
}

// godoc
// @Summary 				  	Get all user instances	
//
// @Schemes
// @Description 				Get all instances for current user ignoring the status	
// @Tags 						instances
//
// @Security					BearerAuth
//
// @Success 					202 {string} string	"accepted"
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /instances [get]
func GetUserInstances(c *gin.Context){
	userId := c.Request.Header.Get(constants.USER_ID_HEADER)

	_,err := db_instances.GetAllUserInstances(userId)

	if err != nil{
		util.ThrowInternalServerErrorException(c,"Could not fetch data")
		return
	}

	c.Status(http.StatusOK)
}