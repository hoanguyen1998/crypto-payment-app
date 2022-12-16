package services

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/hoanguyen1998/crypto-payment-system/helpers"
	"github.com/hoanguyen1998/crypto-payment-system/internal/blockchain"
	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
	"github.com/hoanguyen1998/crypto-payment-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AppService struct {
	repo repository.DatabaseRepo
}

func NewAppService(repo repository.DatabaseRepo) *AppService {
	return &AppService{
		repo: repo,
	}
}

func (s *AppService) NewMasterPublicKey(masterKey string, userId, paymentMethodId int) (int, *helpers.ErrRest) {
	existKey, _ := s.repo.GetMasterKeyByUserIdAndPaymentMethod(userId, paymentMethodId)

	if existKey != nil {
		return 0, helpers.ErrMasterKeyAlreadyExist
	}

	newKey := models.MasterPublicKey{
		PublicKey:       masterKey,
		UserId:          userId,
		NumberOfApps:    0,
		PaymentMethodId: paymentMethodId,
	}

	keyId, err := s.repo.CreateMasterPublicKey(newKey)

	if err != nil {
		return 0, helpers.NewInternalServerError(err.Error())
	}

	return keyId, nil
}

func (s *AppService) NewApp(name string, userId int) (models.Application, *helpers.ErrRest) {
	app := models.Application{
		Name:   name,
		UserId: userId,
	}

	appId, err := s.repo.CreateApp(app)

	if err != nil {
		return app, helpers.NewInternalServerError(err.Error())
	}

	app.Id = appId

	return app, nil
}

func (s *AppService) NewAppKey(appId, paymentMethodId int) (models.ApplicationKey, *helpers.ErrRest) {
	var appKey models.ApplicationKey

	existAppKey, _ := s.repo.GetAppKeyByAppIdAndPaymentMethod(appId, paymentMethodId)

	if existAppKey == nil {
		return appKey, helpers.ErrAppKeyAlreadyExist
	}

	appInfo, err := s.repo.GetAppById(appId)

	if err != nil {
		return appKey, helpers.NewInternalServerError(err.Error())
	}

	existKey, _ := s.repo.GetMasterKeyByUserIdAndPaymentMethod(appInfo.UserId, paymentMethodId)

	if existKey == nil {
		return appKey, helpers.ErrMasterKeyNotFound
	}

	// generate app key
	publicKey, err := blockchain.GenerateApplicationKey(paymentMethodId, existKey.PublicKey, appInfo.AppIndex)

	if err != nil {
		return appKey, helpers.NewInternalServerError(err.Error())
	}

	appKey = models.ApplicationKey{
		PublicKey:         publicKey,
		AppId:             appId,
		MasterPublicKeyId: existKey.Id,
	}

	appKeyId, err := s.repo.CreateAppKey(appKey, existKey.Id)

	if err != nil {
		return appKey, helpers.NewInternalServerError(err.Error())
	}

	appKey.Id = appKeyId

	return appKey, nil
}

func (s *AppService) NewAppAndAppKey(userId, paymentMethodId int, appName string) (models.ApplicationKey, *helpers.ErrRest) {
	var appKey models.ApplicationKey

	existKey, _ := s.repo.GetMasterKeyByUserIdAndPaymentMethod(userId, paymentMethodId)

	if existKey == nil {
		return appKey, helpers.ErrMasterKeyNotFound
	}

	numOfApps := existKey.NumberOfApps + 1

	app := models.Application{
		Name:     appName,
		UserId:   userId,
		AppIndex: numOfApps,
	}

	// generate app key
	publicKey, err := blockchain.GenerateApplicationKey(paymentMethodId, existKey.PublicKey, numOfApps)

	if err != nil {
		return appKey, helpers.NewInternalServerError(err.Error())
	}

	appKey = models.ApplicationKey{
		PublicKey:         publicKey,
		MasterPublicKeyId: existKey.Id,
	}

	appId, appKeyId, err := s.repo.CreateAppAndKey(app, appKey)

	if err != nil {
		return appKey, nil
	}

	appKey.Id = appKeyId
	appKey.AppId = appId

	return appKey, nil
}

func (s *AppService) CreateOrder(appId, paymentMethodId, appOrderId int, amount float64) (models.Order, *helpers.ErrRest) {
	appKey, _ := s.repo.GetAppKeyByAppIdAndPaymentMethod(appId, paymentMethodId)

	var order models.Order

	if appKey == nil {
		return order, helpers.ErrAppKeyNotFound
	}

	addressNum := appKey.AddressGenerated + 1

	appIndex, err := s.repo.GetAppIndexByAppId(appId)

	if err != nil {
		return order, helpers.NewInternalServerError(err.Error())
	}

	path := strconv.Itoa(appIndex) + "/" + strconv.Itoa(addressNum-1)

	// generate address
	address, _ := blockchain.GenerateAddress(paymentMethodId, appKey.PublicKey, addressNum-1)

	order = models.Order{
		ApplicationOrderId: appOrderId,
		Amount:             amount,
		ReceivedAddress:    address,
		Path:               path,
		ApplicationKeyId:   appKey.Id,
	}

	orderId, err := s.repo.CreateOrder(addressNum, appKey.Id, order)

	if err != nil {
		return order, helpers.ErrMasterKeyAlreadyExist
	}

	order.Id = orderId

	return order, nil
}

func (s *AppService) UpdateOrderStatusAndCreateTransaction(orderId int, status string, txData models.Transaction) (models.Transaction, *helpers.ErrRest) {
	txId, err := s.repo.UpdateOrderStatusAndCreateTransaction(orderId, status, txData)

	if err != nil {
		return txData, helpers.NewInternalServerError(err.Error())
	}

	txData.Id = txId

	return txData, nil
}

func (s *AppService) GetOrdersToWithdraw(appKeyId int) ([]models.Order, *helpers.ErrRest) {
	var orders []models.Order

	orders, err := s.repo.GetOrdersToWithdraw(appKeyId)

	if err != nil {
		if err == sql.ErrNoRows {
			return orders, helpers.ErrOrderNotFound
		}

		return orders, helpers.NewInternalServerError(err.Error())
	}

	return orders, nil
}

func (s *AppService) GetLatestBlockNumber(paymentMethodId int) (int, *helpers.ErrRest) {
	blockNumber, err := s.repo.GetLatestBlockNumber(paymentMethodId)

	if err != nil {
		return 0, helpers.NewInternalServerError(err.Error())
	}

	return blockNumber, nil
}

func (s *AppService) CreateUser(email, name, password string) (models.User, *helpers.ErrRest) {
	var storedUser models.User

	userExist, _ := s.repo.GetUserByEmail(email)

	if userExist != nil {
		fmt.Println("exist")
		return storedUser, helpers.ErrEmailAlreadyExist
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return storedUser, helpers.NewInternalServerError(err.Error())
	}

	storedUser = models.User{
		Username:     name,
		PasswordHash: string(passwordHash),
		Email:        email,
	}

	userId, err := s.repo.CreateUser(storedUser)

	if err != nil {
		return storedUser, helpers.NewInternalServerError(err.Error())
	}

	storedUser.Id = userId

	return storedUser, nil
}

func (s *AppService) GetUser(email, password string) (models.User, *helpers.ErrRest) {
	var storedUser models.User

	userExist, _ := s.repo.GetUserByEmail(email)

	if userExist != nil {
		return storedUser, helpers.ErrEmailAlreadyExist
	}

	return *userExist, nil
}
