package repository

import "github.com/hoanguyen1998/crypto-payment-system/internal/models"

type DatabaseRepo interface {
	CreateUser(user models.User) (int, error)
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)

	CreateMasterPublicKey(masterKey models.MasterPublicKey) (int, error)
	GetMasterKeyByUserIdAndPaymentMethod(userId, paymentMethodId int) (*models.MasterPublicKey, error)

	CreateApp(app models.Application) (int, error)
	GetAppById(appId int) (*models.Application, error)
	GetAppsByUserId(userId int) ([]*models.Application, error)
	GetAppIndexByAppId(appId int) (int, error)

	CreateAppKey(appKey models.ApplicationKey, masterPublicKeyId int) (int, error)
	CreateAppAndKey(app models.Application, appKey models.ApplicationKey) (int, int, error)
	GetAppKeyByAppIdAndPaymentMethod(appId int, paymentMethodId int) (*models.ApplicationKey, error)

	CreateOrder(addressNum, appKeyId int, order models.Order) (int, error)
	GetOrdersByAppKey(appKeyId int) ([]*models.Order, error)
	UpdateOrderStatusAndCreateTransaction(orderId int, status string, txData models.Transaction) (int, error)
	GetOrdersToWithdraw(appKeyId int) ([]models.Order, error)

	CreateTransaction(tx models.Transaction) (int, error)
	GetLatestBlockNumber(paymentMethodId int) (int, error)
}
