package domain

type QueueStatus string

const (
	QueueStatusQueued    QueueStatus = "QUEUED"
	QueueStatusOffered   QueueStatus = "OFFERED"
	QueueStatusExpired   QueueStatus = "EXPIRED"
	QueueStatusCompleted QueueStatus = "COMPLETED"
	QueueStatusSoldOut   QueueStatus = "SOLD_OUT"
	QueueStatusCancelled QueueStatus = "CANCELLED"
)
