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
	initializeRuncpaNotifier()
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
	dal := dataAccessLayer{}

	// load balancer status route
	router.GET("/lbstatus", lbstatusHandler)

	// documentation
	router.Static("/doc", "/opt/swagger")
	router.StaticFile("/v1/doc.json", "/opt/apidoc/v1.json")

	// middlewares
	publisherAuthRequiredMiddleware := middlewares.PublisherAuthRequired(dal.GetAuthToken)
	proxyAuthRequiredMiddleware := middlewares.ProxyAuthRequired(dal.GetScoreByIP, config.ProxyDetection.Threshold)
	router.Use(
		middlewares.RecoveryWithWriter(os.Stderr),
		middlewares.Logger(geo),
		middlewares.CORS(),
		gin.ErrorLogger(),
	)

	// integrated with offerwalls
	offerwallHandler := offerwalls.New(dal, runcpaNotifier.CallbackRevenueShare)
	offerwallRouter := router.Group("/offerwalls")

	offerwallRouter.GET("/adgate",
		middlewares.IDParserMiddleware("user_id", models.OfferwallNameAdgate),
		middlewares.AdgateAuthRequired(config.Offerwalls.Adgate.WhitelistIPs),
		offerwallHandler.AdgateCallback(),
	)

	offerwallRouter.GET("/adscend",
		middlewares.IDParserMiddleware("user_id", models.OfferwallNameAdscend),
		middlewares.AdscendAuthRequired(config.Offerwalls.Adscend.WhitelistIPs),
		offerwallHandler.AdscendCallback(),
	)

	offerwallRouter.GET("/kiwiwall",
		middlewares.IDParserMiddleware("sub_id", models.OfferwallNameKiwiwall),
		middlewares.KiwiwallAuthRequired(config.Offerwalls.Kiwiwall.WhitelistIPs, config.Offerwalls.Kiwiwall.SecretKey),
		offerwallHandler.KiwiwallCallback(),
	)

	offerwallRouter.GET("/personaly",
		middlewares.IDParserMiddleware("user_id", models.OfferwallNamePersonaly),
		middlewares.PersonalyAuthRequired(config.Offerwalls.Personaly.WhitelistIPs, config.Offerwalls.Personaly.AppHash, config.Offerwalls.Personaly.SecretKey),
		offerwallHandler.PersonalyCallback(),
	)

	offerwallRouter.GET("/pointclicktrack",
		middlewares.IDParserMiddleware("sid1", models.OfferwallNamePointClickTrack),
		middlewares.PointClickTrackAuthRequired(config.Offerwalls.PointClickTrack.WhitelistIPs),
		offerwallHandler.PointClickTrackCallback(),
	)

	offerwallRouter.GET("/ptcwall",
		middlewares.IDParserMiddleware("usr", models.OfferwallNamePtcwall),
		middlewares.PtcwallAuthRequired(config.Offerwalls.Ptcwall.WhitelistIPs, config.Offerwalls.Ptcwall.PostbackPassword),
		offerwallHandler.PtcwallCallback(),
	)

	offerwallRouter.GET("/wannads",
		middlewares.IDParserMiddleware("subId", models.OfferwallNameWannads),
		middlewares.WannadsAuthRequired(config.Offerwalls.Wannads.SecretKey),
		offerwallHandler.WannadsCallback(),
	)

	// v1 router
	v1Router := router.Group("/v1")

	// v1 user handler
	v1UserHandler := v1.NewUserHandler(dal, dal, runcpaNotifier)
	v1Router.POST("/users", proxyAuthRequiredMiddleware, v1UserHandler.CreateUser())           // create user
	v1Router.GET("/users/:address", proxyAuthRequiredMiddleware, v1UserHandler.RetrieveUser()) // retrieve user

	// v1 publisher handler
	publisherHandler := v1.NewPublisherHandler(dal, dal)
	v1Router.POST("/publishers", publisherHandler.CreatePublisher())                // create publisher account
	v1Router.GET("/publishers/:email", publisherHandler.RetrievePublisherByEmail()) // get publisher info

	// v1 site handler
	v1SiteHandler := v1.NewSiteHandler(dal, dal)
	v1SiteRouter := v1Router.Group("/sites")
	v1SiteRouter.POST("", publisherAuthRequiredMiddleware, v1SiteHandler.CreateSite())   // create site owned by publisher
	v1SiteRouter.GET("", publisherAuthRequiredMiddleware, v1SiteHandler.RetrieveSites()) // get all sites owned by publisher
	v1SiteRouter.GET("/:site_id", v1SiteHandler.RetrieveSite())                          // get one site

	// v1 auth token handler
	authTokenHandler := v1.NewAuthTokenHandler(dal, dal)
	v1Router.POST("/auth/publisher", authTokenHandler.CreateAuthToken()) // create auth token for access to publisher dashboard

	// v1 offers
	v1OfferHandler := v1.NewOfferHandler(dal, dal)
	v1OfferRouter := v1Router.Group("/offers")
	v1OfferRouter.GET("/user/:user_id", proxyAuthRequiredMiddleware, v1OfferHandler.RetrieveOffersByUser())     // get offers filter by user_id
	v1OfferRouter.GET("/site/:site_id", publisherAuthRequiredMiddleware, v1OfferHandler.RetrieveOffersBySite()) // get offers filter by site_id
	v1OfferRouter.GET("/publisher/:publisher_id", v1OfferHandler.RetrieveOffersByPublisher())                   // get offers filter by publisher_id

	// v1 withdrawals
	v1WithdrawalHandler := v1.NewWithdrawalHandler(dal)
	v1WithdrawalRouter := v1Router.Group("/withdrawals")
	v1WithdrawalRouter.GET("/user/:user_id", proxyAuthRequiredMiddleware, v1WithdrawalHandler.RetrieveWithdrawalsByUser())
	v1WithdrawalRouter.GET("/publisher/:publisher_id", v1WithdrawalHandler.RetrieveWithdrawalsByPublisher())

	// start server
	logrus.WithFields(logrus.Fields{
		"event":   models.LogEventChangeServiceState,
		"address": config.HTTP.Address,
		"mode":    config.HTTP.Mode,
		"state":   "running",
	}).Info("server is starting...")
	startServer(config.HTTP.Address, router)
}

func lbstatusHandler(c *gin.Context) {
	status := http.StatusOK
	if !isServiceAlive() {
		status = http.StatusServiceUnavailable
	}
	c.Status(status)
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
