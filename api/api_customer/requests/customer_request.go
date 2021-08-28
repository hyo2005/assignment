package requests

type CreatCustomerRequest struct {
	Name       string `json:"customer_name"`
	Phone      string `json:"phone"`
	License_id string `json:"license_id"`
	Address    string `json:"address"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Active     bool   `json:"active"`
}

type UpdateCustomerRequest struct {
	ID         string `json:"id"`
	Name       string `json:"customer_name"`
	Phone      string `json:"phone"`
	License_id string `json:"license_id"`
	Address    string `json:"address"`
	Email      string `json:"email"`
	Active     bool   `json:"active"`
}
type FindCustomerRequest struct {
	ID string `json:"id" binding:"required"`
}
type ChangePasswordRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
