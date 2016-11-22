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
	offerReader offerReader
	siteReader  siteReader
}

// NewOfferHandler _
func NewOfferHandler(offerReader offerReader, siteReader siteReader) OfferHandler {
	return OfferHandler{
		offerReader: offerReader,
		siteReader:  siteReader,
	}
}

type offerReader interface {
	GetNumberOfOffersByUserID(userID int64) (int64, error)
	GetOffersByUserID(userID, limit, offset int64) ([]rpcmodels.Offer, error)

	GetNumberOfOffersBySiteID(siteID int64) (int64, error)
	GetOffersBySiteID(siteID, limit, offset int64) ([]rpcmodels.Offer, error)

	GetNumberOfOffersByPublisherID(publisherID int64) (int64, error)
	GetOffersByPublisherID(publisherID, limit, offset int64) ([]rpcmodels.Offer, error)
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

		count, err := o.offerReader.GetNumberOfOffersByUserID(userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		offers, err := o.offerReader.GetOffersByUserID(userID, payload.Limit, payload.Offset)
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
		site, _ := o.siteReader.GetSite(siteID)
		authToken := c.MustGet("auth_token").(*rpcmodels.AuthToken)
		if site.PublisherID != authToken.PublisherID {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		count, err := o.offerReader.GetNumberOfOffersBySiteID(siteID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		offers, err := o.offerReader.GetOffersBySiteID(siteID, payload.Limit, payload.Offset)
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

// RetrievePublisherOffers _
func (o OfferHandler) RetrievePublisherOffers() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := paginationPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		publisherID, err := strconv.ParseInt(c.Param("publisher_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		count, err := o.offerReader.GetNumberOfOffersByPublisherID(publisherID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		offers, err := o.offerReader.GetOffersByPublisherID(publisherID, payload.Limit, payload.Offset)
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
