package offerwalls

// // SuperrewardsCallback handles adgate callback
// // NOTE: not ready for production
// func (o OfferwallHandler) SuperrewardsCallback() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		payload := struct {
// 			Amount        int64  `form:"new" binding:"required"`
// 			TransactionID string `form:"id" binding:"required"`
// 			OfferName     string `form:"offer_name"`
// 			Status        int64  `form:"status" binding:"required,eq=1|eq=0"` // 1 success 0 chargeback
// 		}{}
// 		if err := c.BindWith(&payload, binding.Form); err != nil {
// 			return
// 		}

// 		publisherID := c.MustGet("publisher_id").(int64)
// 		siteID := c.MustGet("site_id").(int64)
// 		userID := c.MustGet("user_id").(int64)

// 		offer := models.Offer{
// 			PublisherID:   publisherID,
// 			SiteID:        siteID,
// 			UserID:        userID,
// 			OfferName:     payload.OfferName,
// 			OfferwallName: "superrewards",
// 			TransactionID: payload.TransactionID,
// 			Amount:        payload.Amount,
// 		}

// 		if err := o.handleOfferCallback(offer, payload.Status == 0); err != nil {
// 			c.AbortWithError(http.StatusInternalServerError, err)
// 			return
// 		}

// 		c.String(http.StatusOK, "1")
// 	}
// }
