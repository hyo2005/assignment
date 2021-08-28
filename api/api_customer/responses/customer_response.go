package responses

import "time"

type CustomerResponse struct {
	ID         string `json:"id"`
	Name       string `json:"customer_name"`
	Phone      string `json:"phone"`
	License_id string `json:"license_id"`
	Address    string `json:"address"`
	Email      string `json:"email"`
	Active     bool   `json:"active"`
}

type HistoryResponse struct {
	Customrer_id   string      `json:"customer_id"`
	Booking_code   []string    `json:"code"`
	Booking_id     []string    `json:"booking_id"`
	Flight_id      []string    `json:"flight_id"`
	Booking_status []string    `json:"status"`
	Booking_date   []time.Time `json:"date"`
}
