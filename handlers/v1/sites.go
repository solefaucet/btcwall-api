package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// SiteHandler _
type SiteHandler struct {
	siteReader siteReader
	siteWriter siteWriter
}

// NewSiteHandler _
func NewSiteHandler(siteReader siteReader, siteWriter siteWriter) SiteHandler {
	return SiteHandler{
		siteReader: siteReader,
		siteWriter: siteWriter,
	}
}

type siteWriter interface {
	CreateSite(publisherID int64, siteName, siteURL string) error
}

type siteReader interface {
	GetSite(siteID int64) (*rpcmodels.Site, error)
	GetSitesByPublisherID(publisherID int64) ([]rpcmodels.Site, error)
}

// RetrieveSite _
func (h SiteHandler) RetrieveSite() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteID, err := strconv.ParseInt(c.Param("site_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		site, err := h.siteReader.GetSite(siteID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if site == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, site)
	}
}

// RetrieveSites _
func (h SiteHandler) RetrieveSites() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(*rpcmodels.AuthToken)

		sites, err := h.siteReader.GetSitesByPublisherID(authToken.PublisherID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, sites)
	}
}

// CreateSite _
func (h SiteHandler) CreateSite() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			SiteName string `json:"site_name" binding:"required"`
			SiteURL  string `json:"site_url" binding:"required,url"`
		}{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}

		authToken := c.MustGet("auth_token").(*rpcmodels.AuthToken)
		if err := h.siteWriter.CreateSite(authToken.PublisherID, payload.SiteName, payload.SiteURL); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusCreated)
	}
}
