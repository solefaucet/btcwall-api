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
	withdrawalReader withdrawalReader
}

type withdrawalReader interface {
	GetNumberOfUserWithdrawalsByUserID(userID int64) (int64, error)
	GetUserWithdrawalsByUserID(userID, limit, offset int64) ([]rpcmodels.UserWithdrawal, error)

	GetNumberOfPublisherWithdrawalsByPublisherID(publisherID int64) (int64, error)
	GetPublisherWithdrawalsByPublisherID(publisher, limit, offset int64) ([]rpcmodels.PublisherWithdrawal, error)
}

// NewWithdrawalHandler _
func NewWithdrawalHandler(withdrawalReader withdrawalReader) WithdrawalHandler {
	return WithdrawalHandler{
		withdrawalReader: withdrawalReader,
	}
}

// RetrieveWithdrawalsByUser _
func (h WithdrawalHandler) RetrieveWithdrawalsByUser() gin.HandlerFunc {
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

		count, err := h.withdrawalReader.GetNumberOfUserWithdrawalsByUserID(userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		withdrawals, err := h.withdrawalReader.GetUserWithdrawalsByUserID(userID, payload.Limit, payload.Offset)
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

// RetrieveWithdrawalsByPublisher _
func (h WithdrawalHandler) RetrieveWithdrawalsByPublisher() gin.HandlerFunc {
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

		count, err := h.withdrawalReader.GetNumberOfPublisherWithdrawalsByPublisherID(publisherID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		withdrawals, err := h.withdrawalReader.GetPublisherWithdrawalsByPublisherID(publisherID, payload.Limit, payload.Offset)
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
