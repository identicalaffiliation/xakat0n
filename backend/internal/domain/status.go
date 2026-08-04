package domain

type QueueStatus string

const (
	QueueStatusQueued    QueueStatus = "QUEUED"
	QueueStatusOffered   QueueStatus = "OFFERED"
	QueueStatusCheckout  QueueStatus = "CHECKOUT"
	QueueStatusExpired   QueueStatus = "EXPIRED"
	QueueStatusPurchased QueueStatus = "PURCHASED"
	QueueStatusSoldOut   QueueStatus = "SOLD_OUT"
	QueueStatusCancelled QueueStatus = "CANCELLED"
)
