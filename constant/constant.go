package constant

type OrderStatus string

const (
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusRejected  OrderStatus = "REJECTED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

const (
	CtxKeyRequestID = "request_id"
	CtxKeySubLogger = "sub_logger"
)
