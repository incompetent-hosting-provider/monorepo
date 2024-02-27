package instances

import (
	"errors"
	"incompetent-hosting-provider/backend/pkg/constants"
	db_instances "incompetent-hosting-provider/backend/pkg/db/tables/instances"
	db_presets "incompetent-hosting-provider/backend/pkg/db/tables/presets"
	"incompetent-hosting-provider/backend/pkg/mq_handler"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

type InstanceInfo struct {
	Type               string                    `json:"type"`
	ContainerName      string                    `json:"name"`
	ContainerId        string                    `json:"id"`
	ContainerImageData ContainerImageDescription `json:"container image"`
	InstanceStatus     string                    `json:"status"`
}

type InstanceInfoDetailedResponse struct {
	Type               string                    `json:"type"`
	ContainerName      string                    `json:"name"`
	ContainerId        string                    `json:"id"`
	ContainerImageData ContainerImageDescription `json:"container image"`
	InstanceStatus     string                    `json:"status"`
	StartedAt          string                    `json:"started_at"`
	CreatedAt          string                    `json:"created_at"`
	ContainerPorts     []int                     `json:"open_ports"`
	Description        string                    `json:"description"`
}

type InstancesInfoReponse struct {
	Instances []InstanceInfo `json:"instances"`
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

	preset, err := db_presets.GetPresetById(createRequest.Preset)

	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			util.ThrowNotFoundException(c, "No preset with the provided id exists")
			return
		}
		util.ThrowInternalServerErrorException(c, "Could not delete entry at this time")
		return
	}

	// The env is not persisted as it may contain sensible data
	generatedEnv := map[string]string{}

	for _, v := range preset.RequiredEnv {
		generatedEnv[v] = util.RandStringRunes(64)
	}

	containerId := uuid.NewString()

	// Env with sensible data is currently in eventlog -> This should be fixed via the terraform vault
	err = mq_handler.PublishPresetContainerStartEvent(mq_handler.PresetContainerStartEvent{
		UserId:        userSub,
		ContainerUUID: containerId,
		PresetId:      createRequest.Preset,
		ContainerEnv:  generatedEnv,
	})

	if err != nil {
		util.ThrowServiceUnavailableException(c, "Could not schedule container at the current time")
		return
	}

	err = db_instances.InsertInstance(db_instances.InstancesTable{
		UserSub:              userSub,
		ContainerUUID:        containerId,
		ContainerPorts:       preset.ContainerPorts,
		ContainerDescription: createRequest.Description,
		ContainerName:        createRequest.ContainerName,
		Image:                preset.Image,
		InstanceStatus:       db_instances.STATUS_VALUE_SCHEDULED,
		CreatedAt:            time.Now().Format(time.RFC3339),
		StartedAt:            "N/A",
	})

	if err != nil {
		util.ThrowInternalServerErrorException(c, "Could not save item at the current time")
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

	err = db_instances.InsertInstance(db_instances.InstancesTable{
		UserSub:              userSub,
		ContainerUUID:        containerId,
		ContainerPorts:       createRequest.Ports,
		ContainerDescription: createRequest.Description,
		ContainerName:        createRequest.ContainerName,
		Image: db_instances.ImageSpecification{
			Name: createRequest.Image.ImageName,
			Tag:  createRequest.Image.Tag,
		},
		InstanceStatus: db_instances.STATUS_VALUE_SCHEDULED,
		Type:           db_instances.TYPE_CUSTOM,
		CreatedAt:      time.Now().Format(time.RFC3339),
		StartedAt:      "N/A",
	})

	if err != nil {
		util.ThrowInternalServerErrorException(c, "Could not save item at the current time")
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
// @Router /instances/{containerId} [delete]
func DeleteContainerHandler(c *gin.Context) {
	containerUUID := c.Param("containerId")
	userSub := c.Request.Header.Get(constants.USER_ID_HEADER)

	if containerUUID == "" {
		util.ThrowBadRequestException(c, "No valid containerId passed")
	}

	err := mq_handler.PublishDeleteContainerEvent(mq_handler.DeleteContainerEvent{
		ContainerUUID: containerUUID,
		UserId:        userSub,
	})

	if err != nil {
		util.ThrowServiceUnavailableException(c, "Could  schedule container deletion at the current time")
		return
	}

	err = db_instances.DeleteInstanceById(userSub, containerUUID)
	if err != nil {
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
// @Success 					200 {object} instances.InstancesInfoReponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /instances [get]
func GetUserInstances(c *gin.Context) {
	userSub := c.Request.Header.Get(constants.USER_ID_HEADER)

	instances, err := db_instances.GetAllUserInstances(userSub)

	if err != nil {
		util.ThrowInternalServerErrorException(c, "Could not fetch data")
		return
	}

	c.JSON(http.StatusOK, InstancesInfoReponse{
		Instances: serializeInstanceResponses(instances),
	})
}

// godoc
// @Summary 				  	Get instance details
//
// @Schemes
// @Description 				Get details of a single instance by id
// @Tags 						instances
//
// @Security					BearerAuth
//
// @Param   containerId     path    string     true        "Container Id"
//
// @Success 					200 {object} instances.InstanceInfoDetailedResponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /instances/{containerId} [get]
func GetInstance(c *gin.Context) {
	containerUUID := c.Param("containerId")
	userSub := c.Request.Header.Get(constants.USER_ID_HEADER)

	if containerUUID == "" {
		util.ThrowBadRequestException(c, "No valid containerId passed")
	}

	instance, err := db_instances.GetInstanceById(userSub, containerUUID)

	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			util.ThrowNotFoundException(c, "Could not find item with given id")
			return
		}
		util.ThrowInternalServerErrorException(c, "Could not delete entry at this time")
		return
	}

	c.JSON(http.StatusOK, serializeDetailedInstanceResponse(instance))
}
