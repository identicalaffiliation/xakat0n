package dto

// CheckoutStarted — ответ POST /items/{itemId}/checkout.
type CheckoutStarted struct {
	QueueApplied bool    `json:"queueApplied"`
	Ticket       *Ticket `json:"ticket"`
}

func NewCheckoutStarted(ticket *Ticket) *CheckoutStarted {
	return &CheckoutStarted{
		QueueApplied: ticket != nil,
		Ticket:       ticket,
	}
}
