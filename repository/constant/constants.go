package constant

import myutils "gwlkm-resend-transaction/myutlis"

// Convert Mysql Data Type
var (
	SQLsaldo_before_trans          myutils.FieldInt
	SQLsynced_ibs_core_description myutils.FieldString
	SQLid_user                     myutils.FieldInt
	SQLid_raw                      myutils.FieldInt
	SQLno_hp_alternatif            myutils.FieldString
)

var (
	Success = "0000"
	Pending = "1234"
	Failed  = "1100"
)
