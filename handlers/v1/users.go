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
	userStorage userStorage
}

// NewUserHandler _
func NewUserHandler(userStorage userStorage) UserHandler {
	return UserHandler{
		userStorage: userStorage,
	}
}

type userStorage interface {
	userReader
	userWriter
}

type userReader interface {
	GetUser(address string) (*rpcmodels.User, error)
}

type userWriter interface {
	CreateUser(address string) (*rpcmodels.User, error)
}

// CreateUser creates user, response with user info
func (userHandler UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Address string `json:"address" binding:"required"`
		}{}

		if err := c.BindJSON(&payload); err != nil {
			return
		}
		payload.Address = strings.TrimSpace(payload.Address)
		if !validateBitcoinAddress(payload.Address) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bitcoin address %s is invalid", payload.Address))
			return
		}

		user, err := userHandler.userStorage.CreateUser(payload.Address)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if user == nil {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// RetrieveUser response with user info
func (userHandler UserHandler) RetrieveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		user, err := userHandler.userStorage.GetUser(address)
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
