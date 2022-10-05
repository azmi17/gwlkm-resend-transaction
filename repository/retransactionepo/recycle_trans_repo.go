package retransactionepo

import "gwlkm-resend-transaction/entities"

type RetransactionRepo interface {
	RecycleTransaction(dataTrans *entities.MsgTransHistory) error
	RecycleReversedTransaction(dataTrans *entities.MsgTransHistory) error
}
