package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// WithdrawalHandler _
type WithdrawalHandler struct {
	dependency withdrawHandlerDependency
}

type withdrawHandlerDependency interface {
	withdrawalReader
}

type withdrawalReader interface {
	GetNumberOfUserWithdrawalsByUserID(userID int64) (int64, error)
	GetUserWithdrawalsByUserID(userID, limit, offset int64) ([]rpcmodels.UserWithdrawal, error)

	GetNumberOfPublisherWithdrawalsByPublisherID(publisherID int64) (int64, error)
	GetPublisherWithdrawalsByPublisherID(publisher, limit, offset int64) ([]rpcmodels.PublisherWithdrawal, error)
}

// NewWithdrawalHandler _
func NewWithdrawalHandler(dependency withdrawHandlerDependency) WithdrawalHandler {
	return WithdrawalHandler{
		dependency: dependency,
	}
}

// UserWithdrawalHandler _
func (w WithdrawalHandler) UserWithdrawalHandler() gin.HandlerFunc {
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

		count, err := w.dependency.GetNumberOfUserWithdrawalsByUserID(userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		withdrawals, err := w.dependency.GetUserWithdrawalsByUserID(userID, payload.Limit, payload.Offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, paginationResult{
			Count: count,
			Data:  withdrawals,
		})
	}
}

// PublisherWithdrawalHandler _
func (w WithdrawalHandler) PublisherWithdrawalHandler() gin.HandlerFunc {
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

		count, err := w.dependency.GetNumberOfPublisherWithdrawalsByPublisherID(publisherID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		withdrawals, err := w.dependency.GetPublisherWithdrawalsByPublisherID(publisherID, payload.Limit, payload.Offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, paginationResult{
			Count: count,
			Data:  withdrawals,
		})
	}
}
