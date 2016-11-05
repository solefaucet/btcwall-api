package v1

import (
	"net/http"
	"strconv"

	rpcmodels "github.com/solefaucet/btcwall-rpc-model"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// OfferHandler _
type OfferHandler struct {
	dependency offerHandlerDependency
}

type offerHandlerDependency interface {
	offerReader
	siteReader
}

// NewOfferHandler _
func NewOfferHandler(dependency offerHandlerDependency) OfferHandler {
	return OfferHandler{
		dependency: dependency,
	}
}

type offerReader interface {
	GetNumberOfOffersByUserID(userID int64) (int64, error)
	GetOffersByUserID(userID, limit, offset int64) ([]rpcmodels.Offer, error)

	GetNumberOfOffersBySiteID(siteID int64) (int64, error)
	GetOffersBySiteID(siteID, limit, offset int64) ([]rpcmodels.Offer, error)
}

// UserOfferHandler _
func (o OfferHandler) UserOfferHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := paginationPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		count, err := o.dependency.GetNumberOfOffersByUserID(userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		offers, err := o.dependency.GetOffersByUserID(userID, payload.Limit, payload.Offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, paginationResult{
			Count: count,
			Data:  offers,
		})
	}
}

// SiteOfferHandler _
func (o OfferHandler) SiteOfferHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := paginationPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		siteID, err := strconv.ParseInt(c.Param("site_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// check if authorized
		site, _ := o.dependency.GetSite(siteID)
		authToken := c.MustGet("auth_token").(*rpcmodels.AuthToken)
		if site.PublisherID != authToken.PublisherID {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		count, err := o.dependency.GetNumberOfOffersBySiteID(siteID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		offers, err := o.dependency.GetOffersBySiteID(siteID, payload.Limit, payload.Offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, paginationResult{
			Count: count,
			Data:  offers,
		})
	}
}
