package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/solefaucet/btcwall-api/handlers/offerwalls"
	"github.com/solefaucet/btcwall-api/handlers/v1"
	"github.com/solefaucet/btcwall-api/middlewares"
	"github.com/solefaucet/btcwall-api/models"
	"github.com/spf13/viper"
)

var (
	geo *geoip2.Reader
)

func main() {
	// initialization
	initializeConfiguration()
	initializeLogging()
	initializeRPCClient()
	geo = must(geoip2.Open(config.Geo.Filename)).(*geoip2.Reader)
	log.Println(spew.Sdump(config))

	// watch configuration changed
	viper.WatchConfig()
	viper.OnConfigChange(func(fsnotify.Event) {
		var c models.Configuration
		if err := viper.Unmarshal(&c); err != nil {
			logrus.WithFields(logrus.Fields{
				"event": models.LogEventReloadConfiguration,
				"error": err.Error(),
			}).Error("fail to reload configuration")
			return
		}
		if err := c.Validate(); err != nil {
			logrus.WithFields(logrus.Fields{
				"event": models.LogEventReloadConfiguration,
				"error": err.Error(),
			}).Error("fail to validate configuration")
			return
		}

		// NOTE: need a mutex to protect config if there is concurrent access
		config.Copy(c)
		log.Println(spew.Sdump(config))

		// reload logging
		initializeLogging()
	})

	// http server
	router := gin.New()
	var dal dataAccessLayer

	// middlewares
	publisherAuthRequiredMiddleware := middlewares.PublisherAuthRequired(dal.GetAuthToken)
	router.Use(
		middlewares.RecoveryWithWriter(os.Stderr),
		middlewares.Logger(geo),
		middlewares.CORS(),
		gin.ErrorLogger(),
	)

	// load balancer status route
	router.GET("/lbstatus", func(c *gin.Context) {
		status := http.StatusOK
		if !isServiceAlive() {
			status = http.StatusServiceUnavailable
		}
		c.Status(status)
	})

	// integrated with offerwalls
	offerwallHandler := offerwalls.New(dal)
	offerwallRouter := router.Group("/offerwalls")

	offerwallRouter.GET("/adgate",
		middlewares.IDParserMiddleware("user_id"),
		middlewares.AdgateAuthRequired(config.Offerwalls.Adgate.WhitelistIPs),
		offerwallHandler.AdgateCallback(),
	)

	offerwallRouter.GET("/adscend",
		middlewares.IDParserMiddleware("user_id"),
		middlewares.AdscendAuthRequired(config.Offerwalls.Adscend.WhitelistIPs),
		offerwallHandler.AdscendCallback(),
	)

	offerwallRouter.GET("/kiwiwall",
		middlewares.IDParserMiddleware("sub_id"),
		middlewares.KiwiwallAuthRequired(config.Offerwalls.Kiwiwall.WhitelistIPs, config.Offerwalls.Kiwiwall.SecretKey),
		offerwallHandler.KiwiwallCallback(),
	)

	offerwallRouter.GET("/personaly",
		middlewares.IDParserMiddleware("user_id"),
		middlewares.PersonalyAuthRequired(config.Offerwalls.Personaly.WhitelistIPs, config.Offerwalls.Personaly.AppHash, config.Offerwalls.Personaly.SecretKey),
		offerwallHandler.PersonalyCallback(),
	)

	offerwallRouter.GET("/pointclicktrack",
		middlewares.IDParserMiddleware("sid1"),
		middlewares.PointClickTrackAuthRequired(config.Offerwalls.PointClickTrack.WhitelistIPs),
		offerwallHandler.PointClickTrackCallback(),
	)

	offerwallRouter.GET("/ptcwall",
		middlewares.IDParserMiddleware("usr"),
		middlewares.PtcwallAuthRequired(config.Offerwalls.Ptcwall.WhitelistIPs, config.Offerwalls.Ptcwall.PostbackPassword),
		offerwallHandler.PtcwallCallback(),
	)

	offerwallRouter.GET("/wannads",
		middlewares.IDParserMiddleware("subId"),
		middlewares.WannadsAuthRequired(config.Offerwalls.Wannads.SecretKey),
		offerwallHandler.WannadsCallback(),
	)

	// v1 router
	v1Router := router.Group("/v1")

	// v1 user handler
	v1UserHandler := v1.NewUserHandler(dal)
	v1Router.POST("/users", v1UserHandler.CreateUser())           // create user
	v1Router.GET("/users/:address", v1UserHandler.RetrieveUser()) // retrieve user

	// v1 publisher handler
	publisherHandler := v1.NewPublisherHandler(dal)
	v1Router.POST("/publishers", publisherHandler.CreatePublisher())                                                 // create publisher account
	v1Router.GET("/publishers/:publisher_id", publisherAuthRequiredMiddleware, publisherHandler.RetrievePublisher()) // check publisher info

	// v1 site handler
	v1SiteHandler := v1.NewSiteHandler(dal)
	v1SiteRouter := v1Router.Group("/sites")
	v1SiteRouter.POST("", publisherAuthRequiredMiddleware, v1SiteHandler.CreateSite())   // create site owned by publisher
	v1SiteRouter.GET("", publisherAuthRequiredMiddleware, v1SiteHandler.RetrieveSites()) // get all sites owned by publisher
	v1SiteRouter.GET("/:site_id", v1SiteHandler.RetrieveSite())                          // get one site

	// v1 auth token handler
	v1Router.POST("/auth/publisher", v1.CreateAuthToken(dal.GetPublisher, dal.CreateAuthToken)) // create auth token for access to publisher dashboard

	// v1 offers
	v1OfferHandler := v1.NewOfferHandler(dal)
	v1OfferRouter := v1Router.Group("/offers")
	v1OfferRouter.GET("/user/:user_id", v1OfferHandler.UserOfferHandler())                                  // get offers filter by user_id
	v1OfferRouter.GET("/site/:site_id", publisherAuthRequiredMiddleware, v1OfferHandler.SiteOfferHandler()) // get offers filter by site_id

	// start server
	logrus.WithFields(logrus.Fields{
		"event":   models.LogEventChangeServiceState,
		"address": config.HTTP.Address,
		"mode":    config.HTTP.Mode,
		"state":   "running",
	})
	startServer(config.HTTP.Address, router)
}

func startServer(address string, handler http.Handler) {
	server := endless.NewServer(address, handler)
	for _, signal := range signals {
		server.RegisterSignalHook(endless.PRE_SIGNAL, signal, preSignalHook)
		server.RegisterSignalHook(endless.POST_SIGNAL, signal, postSignalHook)
	}

	if err := server.ListenAndServe(); err != nil {
		if isServiceAlive() {
			logrus.Errorf("fail to start http server: %v", err)
			return
		}

		logrus.Info("server is shutdown gracefully")
	}
}

func preSignalHook() {
	logrus.WithFields(logrus.Fields{
		"event": models.LogEventChangeServiceState,
		"state": "stopping",
	}).Info("service is stopping...")

	setServiceToDead()
}

func postSignalHook() {
	logrus.WithFields(logrus.Fields{
		"event": models.LogEventChangeServiceState,
		"state": "stopped",
	}).Info("service is stopped...")
}

// helpers
//
// fail fast on initialization
func must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return i
}
