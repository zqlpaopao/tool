package src

type Flags uint

const (
	_ Flags = iota
	RequestGet
	RequestPost
)

var flagNames = []string{
	"GET",
	"POST",
}

type RequestArgs struct {
	RequestUrl string
	RequestParams map[string]string
	RequestType Flags
	RequestHeader map[string]string
	SetTimeOut int64
}

func (f Flags) String() string {
	return  flagNames[f-1]
}

func (f Flags)Int()int{
	return int(f)
}

//ContentTypeFlags ----------------------------------- request header -----------------------------------//
type ContentTypeFlags uint

const (
	_ ContentTypeFlags  = iota
	RequestForm
	RequestRaw
)

var ContentTypeFlagsNames = []map[string]string{
	{"Content-Type": "application/x-www-form-urlencoded"},
	{"Content-Type": "application/json"},
}


func (f ContentTypeFlags) String() map[string]string {
	return  ContentTypeFlagsNames[f-1]
}

func (f ContentTypeFlags)Int()int{
	return int(f)
}