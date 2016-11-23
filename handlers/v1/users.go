package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// UserHandler _
type UserHandler struct {
	userReader                 userReader
	userWriter                 userWriter
	runcpaRegistrationNotifier runcpaRegistrationNotifier
}

// NewUserHandler _
func NewUserHandler(
	userReader userReader, userWriter userWriter,
	runcpaRegistrationNotifier runcpaRegistrationNotifier,
) UserHandler {
	return UserHandler{
		userReader:                 userReader,
		userWriter:                 userWriter,
		runcpaRegistrationNotifier: runcpaRegistrationNotifier,
	}
}

type userReader interface {
	GetUser(address string) (*rpcmodels.User, error)
}

type userWriter interface {
	CreateUser(address, trackID string) (*rpcmodels.User, error)
}

type runcpaRegistrationNotifier interface {
	CallbackRegistration(trackID string)
}

// CreateUser creates user, response with user info
func (userHandler UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Address string `json:"address" binding:"required"`
			TrackID string `json:"track_id"`
		}{}

		if err := c.BindJSON(&payload); err != nil {
			return
		}
		logrus.WithFields(logrus.Fields{
			"event":    models.LogEventCreateUser,
			"address":  payload.Address,
			"track_id": payload.TrackID,
		}).Debug("create user")

		payload.Address = strings.TrimSpace(payload.Address)
		if !validateBitcoinAddress(payload.Address) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bitcoin address %s is invalid", payload.Address))
			return
		}

		user, err := userHandler.userWriter.CreateUser(payload.Address, payload.TrackID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if user == nil {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		// callback to runcpa
		userHandler.runcpaRegistrationNotifier.CallbackRegistration(payload.TrackID)

		c.JSON(http.StatusCreated, user)
	}
}

// RetrieveUser response with user info
func (userHandler UserHandler) RetrieveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		user, err := userHandler.userReader.GetUser(address)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if user == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
