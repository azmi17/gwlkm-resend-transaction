package retransactionepo

import "gwlkm-resend-transaction/entities"

type RetransactionRepo interface {
	RecycleTransaction(dataTrans *entities.IsoMessageBody) error
	RecycleGwlkmTransaction(dataTrans *entities.IsoMessageBody) error
	RecycleLkmTransferSMprematureRevOnCre(dataTrans *entities.IsoMessageBody) error
	ResendReversalGwlkmTransaction(dataTrans *entities.IsoMessageBody) error
	ResendReversalBeforeRecycleGwlkmTransaction(stan string) error
}
