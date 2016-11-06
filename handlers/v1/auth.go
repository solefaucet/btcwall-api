package v1

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
	"github.com/twinj/uuid"
)

// AuthTokenHandler _
type AuthTokenHandler struct {
	publisherReader publisherReader
	authTokenWriter authTokenWriter
}

type authTokenWriter interface {
	CreateAuthToken(authToken rpcmodels.AuthToken) error
}

// NewAuthTokenHandler _
func NewAuthTokenHandler(publisherReader publisherReader, authTokenWriter authTokenWriter) AuthTokenHandler {
	return AuthTokenHandler{
		publisherReader: publisherReader,
		authTokenWriter: authTokenWriter,
	}
}

// CreateAuthToken create auth token for publisher
func (h AuthTokenHandler) CreateAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}

		publisher, err := h.publisherReader.GetPublisher(payload.Email)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if publisher == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(publisher.Password), []byte(payload.Password)); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := rpcmodels.AuthToken{
			PublisherID: publisher.ID,
			Email:       publisher.Email,
			AuthToken:   uuid.NewV4().String(),
		}
		if err := h.authTokenWriter.CreateAuthToken(authToken); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, authToken)
	}
}
