package retransactionepo

import "gwlkm-resend-transaction/entities"

type RetransactionRepo interface {
	RecycleTransaction(dataTrans *entities.IsoMessageBody) error
	RecycleGwlkmTransaction(dataTrans *entities.IsoMessageBody) error
}
