package entities

const (
	PRINTOUT_TYPE_LOG = iota
	PRINTOUT_TYPE_ERR
)

var (
	PrintOutChan = make(chan PrintOut)
)

type PrintOut struct {
	Type    int
	Message []interface{}
}

func PrintError(message ...interface{}) {
	po := PrintOut{
		Type:    PRINTOUT_TYPE_ERR,
		Message: message,
	}

	PrintOutChan <- po
}

func PrintLog(message ...interface{}) {
	po := PrintOut{
		Type:    PRINTOUT_TYPE_LOG,
		Message: message,
	}
	PrintOutChan <- po
}
