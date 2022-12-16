package models

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

type PaymentMethod struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MasterPublicKey struct {
	Id              int    `json:"id"`
	PublicKey       string `json:"public_key"`
	UserId          int    `json:"user_id"`
	PaymentMethodId int    `json:"payment_method_id"`
	NumberOfApps    int    `json:"number_of_apps"`
}

type Application struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	AppIndex int    `json:"app_index"`
	UserId   int    `json:"user_id"`
}

type ApplicationKey struct {
	Id                int    `json:"id"`
	PublicKey         string `json:"public_key"`
	AddressGenerated  int    `json:"address_generated"`
	AppId             int    `json:"app_id"`
	MasterPublicKeyId int    `json:"master_public_key_id"`
}

type Order struct {
	Id                 int     `json:"id"`
	ApplicationOrderId int     `json:"app_order_id"`
	Amount             float64 `json:"amount"`
	ReceivedAddress    string  `json:"received_address"`
	Status             string  `json:"status"`
	Path               string  `json:"path"`
	ApplicationKeyId   int     `json:"application_key_id"`
}

type Transaction struct {
	Id              int    `json:"id"`
	TxHash          string `json:"tx_hash"`
	Sender          string `json:"sender"`
	Recipient       string `json:"recipient"`
	Amount          int    `json:"amount"`
	PaymentMethodId int    `json:"payment_method_id"`
	BlockNumber     int    `json:"block_number"`
	OrderId         int    `json:"order_id"`
	TxType          string `json:"tx_type"`
}
