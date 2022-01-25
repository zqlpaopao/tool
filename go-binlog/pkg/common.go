package pkg

const (
	MinProtocolVersion byte = 10

	OK_HEADER          byte = 0x00
	ERR_HEADER         byte = 0xff
	EOF_HEADER         byte = 0xfe
	LocalInFile_HEADER byte = 0xfb
)
type Config struct {
	Host string
	Port int
	User string
	Pass string
	ServerId int
	LogFile string
	Position int
}
const (
	MaxPayloadLength = 1<<24 - 1
	MinPort = 1024
	ConnTime = 10
	ReadFrom = "last event read from"
)


const (
	SuccessCode = iota+1
	FailCode
)

type errInfo struct {
	code uint16
	state string
	msg string
}



type dataInfo struct{
	code int
	errInfo *errInfo
}

