package dto

import "github.com/google/uuid"

type PaymentCallbackRequest struct {
	TicketID uuid.UUID `json:"ticketId"`
	Result   string    `json:"result"`
}

func (r *PaymentCallbackRequest) Paid() bool {
	return r.Result == "paid"
}

func (r *PaymentCallbackRequest) Valid() bool {
	return r.TicketID != uuid.Nil && (r.Result == "paid" || r.Result == "failed")
}
