package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// SiteHandler _
type SiteHandler struct {
	dependency siteHandlerDependency
}

// NewSiteHandler _
func NewSiteHandler(dependency siteHandlerDependency) SiteHandler {
	return SiteHandler{
		dependency: dependency,
	}
}

type siteHandlerDependency interface {
	siteReader
	siteWriter
}

type siteWriter interface {
	CreateSite(publisherID int64, siteName, siteURL string) error
}

type siteReader interface {
	GetSite(siteID int64) (*rpcmodels.Site, error)
	GetSitesByPublisherID(publisherID int64) ([]rpcmodels.Site, error)
}

// RetrieveSite _
func (s SiteHandler) RetrieveSite() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteID, err := strconv.ParseInt(c.Param("site_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		site, err := s.dependency.GetSite(siteID)
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
func (s SiteHandler) RetrieveSites() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(*rpcmodels.AuthToken)

		sites, err := s.dependency.GetSitesByPublisherID(authToken.PublisherID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, sites)
	}
}

// CreateSite _
func (s SiteHandler) CreateSite() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			SiteName string `json:"site_name" binding:"required"`
			SiteURL  string `json:"site_url" binding:"required,url"`
		}{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}

		authToken := c.MustGet("auth_token").(*rpcmodels.AuthToken)
		if err := s.dependency.CreateSite(authToken.PublisherID, payload.SiteName, payload.SiteURL); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusCreated)
	}
}
