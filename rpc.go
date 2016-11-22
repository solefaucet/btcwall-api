package main

import (
	"github.com/hprose/hprose-golang/rpc"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

var remote *stub

func initializeRPCClient() {
	client := rpc.NewHTTPClient(config.RPC.Address)
	client.MaxIdleConnsPerHost = config.RPC.MaxIdleConnsPerHost
	client.MaxConcurrentRequests = config.RPC.MaxConcurrentRequests
	client.UseService(&remote)
}

type stub struct {
	// auth token
	GetAuthToken    func(authToken string) (*rpcmodels.AuthToken, error)
	CreateAuthToken func(authToken rpcmodels.AuthToken) error

	// offer
	GetNumberOfOffersByUserID func(userID int64) (int64, error)
	GetNumberOfOffersBySiteID func(siteID int64) (int64, error)
	GetOffersByUserID         func(userID, limit, offset int64) ([]rpcmodels.Offer, error)
	GetOffersBySiteID         func(siteID, limit, offset int64) ([]rpcmodels.Offer, error)
	CreateOffer               func(offer rpcmodels.Offer) (duplicated bool, _ error)
	ChargebackOffer           func(offer rpcmodels.Offer) (alreadyChargeback bool, _ error)

	GetNumberOfOffersByPublisherID func(publisherID int64) (int64, error)
	GetOffersByPublisherID         func(publisherID, limit, offset int64) ([]rpcmodels.Offer, error)

	// publisher
	GetPublisher    func(email string) (*rpcmodels.Publisher, error)
	CreatePublisher func(email, password, address string) (*rpcmodels.Publisher, error)

	// site
	GetSite               func(siteID int64) (*rpcmodels.Site, error)
	GetSitesByPublisherID func(publisherID int64) ([]rpcmodels.Site, error)
	CreateSite            func(publisherID int64, siteName, siteURL string) error

	// user
	GetUser    func(address string) (*rpcmodels.User, error)
	CreateUser func(address, trackID string) (*rpcmodels.User, error)

	// withdrawals
	GetNumberOfUserWithdrawalsByUserID           func(userID int64) (int64, error)
	GetUserWithdrawalsByUserID                   func(userID, limit, offset int64) ([]rpcmodels.UserWithdrawal, error)
	GetNumberOfPublisherWithdrawalsByPublisherID func(publisherID int64) (int64, error)
	GetPublisherWithdrawalsByPublisherID         func(publisherID, limit, offset int64) ([]rpcmodels.PublisherWithdrawal, error)

	// proxy detection
	GetScoreByIP func(ip string) (int64, error)
}

type dataAccessLayer struct{}

func (dataAccessLayer) GetAuthToken(authToken string) (*rpcmodels.AuthToken, error) {
	return remote.GetAuthToken(authToken)
}

func (dataAccessLayer) CreateAuthToken(authToken rpcmodels.AuthToken) error {
	return remote.CreateAuthToken(authToken)
}

func (dataAccessLayer) GetNumberOfOffersByUserID(userID int64) (int64, error) {
	return remote.GetNumberOfOffersByUserID(userID)
}

func (dataAccessLayer) GetNumberOfOffersBySiteID(siteID int64) (int64, error) {
	return remote.GetNumberOfOffersBySiteID(siteID)
}

func (dataAccessLayer) GetNumberOfOffersByPublisherID(publisherID int64) (int64, error) {
	return remote.GetNumberOfOffersByPublisherID(publisherID)
}

func (dataAccessLayer) GetOffersByUserID(userID, limit, offset int64) ([]rpcmodels.Offer, error) {
	return remote.GetOffersByUserID(userID, limit, offset)
}

func (dataAccessLayer) GetOffersBySiteID(siteID, limit, offset int64) ([]rpcmodels.Offer, error) {
	return remote.GetOffersBySiteID(siteID, limit, offset)
}

func (dataAccessLayer) GetOffersByPublisherID(publisherID, limit, offset int64) ([]rpcmodels.Offer, error) {
	return remote.GetOffersByPublisherID(publisherID, limit, offset)
}

func (dataAccessLayer) CreateOffer(offer rpcmodels.Offer) (duplicated bool, _ error) {
	return remote.CreateOffer(offer)
}

func (dataAccessLayer) ChargebackOffer(offer rpcmodels.Offer) (alreadyChargeback bool, _ error) {
	return remote.ChargebackOffer(offer)
}

func (dataAccessLayer) GetPublisher(email string) (*rpcmodels.Publisher, error) {
	return remote.GetPublisher(email)
}

func (dataAccessLayer) CreatePublisher(email, password, address string) (*rpcmodels.Publisher, error) {
	return remote.CreatePublisher(email, password, address)
}

func (dataAccessLayer) GetSite(siteID int64) (*rpcmodels.Site, error) {
	return remote.GetSite(siteID)
}

func (dataAccessLayer) GetSitesByPublisherID(publisherID int64) ([]rpcmodels.Site, error) {
	return remote.GetSitesByPublisherID(publisherID)
}

func (dataAccessLayer) CreateSite(publisherID int64, siteName, siteURL string) error {
	return remote.CreateSite(publisherID, siteName, siteURL)
}

func (dataAccessLayer) GetUser(address string) (*rpcmodels.User, error) {
	return remote.GetUser(address)
}

func (dataAccessLayer) CreateUser(address, trackID string) (*rpcmodels.User, error) {
	return remote.CreateUser(address, trackID)
}

func (dataAccessLayer) GetUserWithdrawalsByUserID(userID, limit, offset int64) ([]rpcmodels.UserWithdrawal, error) {
	return remote.GetUserWithdrawalsByUserID(userID, limit, offset)
}

func (dataAccessLayer) GetPublisherWithdrawalsByPublisherID(publisherID, limit, offset int64) ([]rpcmodels.PublisherWithdrawal, error) {
	return remote.GetPublisherWithdrawalsByPublisherID(publisherID, limit, offset)
}

func (dataAccessLayer) GetNumberOfUserWithdrawalsByUserID(userID int64) (int64, error) {
	return remote.GetNumberOfUserWithdrawalsByUserID(userID)
}

func (dataAccessLayer) GetNumberOfPublisherWithdrawalsByPublisherID(publisherID int64) (int64, error) {
	return remote.GetNumberOfPublisherWithdrawalsByPublisherID(publisherID)
}

func (dataAccessLayer) GetScoreByIP(ip string) (int64, error) {
	return remote.GetScoreByIP(ip)
}
