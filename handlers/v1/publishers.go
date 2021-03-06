package v1

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// PublisherHandler _
type PublisherHandler struct {
	publisherReader publisherReader
	publisherWriter publisherWriter
}

// NewPublisherHandler _
func NewPublisherHandler(publisherReader publisherReader, publisherWriter publisherWriter) PublisherHandler {
	return PublisherHandler{
		publisherReader: publisherReader,
		publisherWriter: publisherWriter,
	}
}

type publisherReader interface {
	GetPublisher(email string) (*rpcmodels.Publisher, error)
}

type publisherWriter interface {
	CreatePublisher(email, password, address string) (*rpcmodels.Publisher, error)
}

// RetrievePublisher response with publisher info
func (h PublisherHandler) RetrievePublisher() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(*rpcmodels.AuthToken)

		publisher, err := h.publisherReader.GetPublisher(authToken.Email)

		if err != nil || publisher == nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if fmt.Sprint(publisher.ID) != c.Param("publisher_id") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.JSON(http.StatusOK, publisher)
	}
}

// RetrievePublisherByEmail response with publisher info
func (h PublisherHandler) RetrievePublisherByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		publisher, err := h.publisherReader.GetPublisher(email)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if publisher == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, publisher)
	}
}

// CreatePublisher creates publisher
func (h PublisherHandler) CreatePublisher() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
			Address  string `json:"address" binding:"required,btc_addr"`
		}{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}

		// hash password
		password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		payload.Password = string(password)

		// create publisher
		publisher, err := h.publisherWriter.CreatePublisher(payload.Email, payload.Password, payload.Address)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if publisher == nil {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.JSON(http.StatusCreated, publisher)
	}
}
