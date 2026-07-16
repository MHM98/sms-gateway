package controllermodel

type CreateMessage struct {
	UserID       uint64
	ChargeAmount uint64
	Recipient    string
	Body         string
	ServiceType  string
	Status       string
}
