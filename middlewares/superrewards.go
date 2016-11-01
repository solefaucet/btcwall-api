package middlewares

// // SuperrewardsAuthRequired rejects request if client ip is not in the list
// func SuperrewardsAuthRequired(whitelistIPs string, secretKey string) gin.HandlerFunc {
// 	ips := make(map[string]struct{})
// 	for _, v := range strings.Split(whitelistIPs, ",") {
// 		ips[v] = struct{}{}
// 	}

// 	return func(c *gin.Context) {
// 		if _, ok := ips[c.ClientIP()]; !ok {
// 			c.AbortWithStatus(http.StatusForbidden)
// 			return
// 		}

// 		data := fmt.Sprintf("%v:%v:%v:%v", c.Query("id"), c.Query("new"), c.Query("uid"), secretKey)
// 		if sign := fmt.Sprintf("%x", md5.Sum([]byte(data))); sign != c.Query("sig") {
// 			httprequest, _ := httputil.DumpRequest(c.Request, true)
// 			logrus.WithFields(logrus.Fields{
// 				"event":          models.LogEventKiwiwallInvalidSignature,
// 				"signature":      sign,
// 				"q_signature":    c.Query("sig"),
// 				"request":        string(httprequest),
// 				"id_combination": c.Query("uid"),
// 			}).Debug("signature not match")
// 			return
// 		}

// 		c.Set("publisher_id", publisherID)
// 		c.Set("site_id", siteID)
// 		c.Set("user_id", userID)

// 		c.Next()
// 	}
// }
