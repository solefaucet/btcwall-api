package offerwalls

// // OffertoroCallback handles offertoro callback
// func (o OfferwallHandler) OffertoroCallback() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		payload := struct {
// 			Amount        int64  `form:"amount" binding:"required"`
// 			TransactionID string `form:"id" binding:"required"`
// 			OfferName     string `form:"o_name"`
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
// 			OfferwallName: "offertoro",
// 			TransactionID: payload.TransactionID,
// 			Amount:        payload.Amount,
// 		}

// 		// NOTE: no chargeback for offertoro
// 		if err := o.handleOfferCallback(offer, false); err != nil {
// 			c.AbortWithError(http.StatusInternalServerError, err)
// 			return
// 		}

// 		c.String(http.StatusOK, "1")
// 	}
// }
