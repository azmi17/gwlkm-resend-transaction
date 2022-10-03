package retransactionepo

import "gwlkm-resend-transaction/entities"

type RetransactionRepo interface {
	RecycleTransaction(dataTrans entities.TransHisotry) (string, error)
}
