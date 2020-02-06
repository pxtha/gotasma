package status

type (
	Status struct {
		XCode    uint32 `json:"code" yaml:"code"`
		XStatus  int    `json:"status" yaml:"status"`
		XMessage string `json:"message" yaml:"message"`
	}
)

// New return a new status.
func New(code uint32, status int, message string) Status {
	return Status{
		XCode:    code,
		XStatus:  status,
		XMessage: message,
	}
}

func (s Status) Error() string {
	return s.XMessage
}

func (s Status) Code() uint32 {
	return s.XCode
}

func (s Status) Message() string {
	return s.XMessage
}

func (s Status) Status() int {
	return s.XStatus
}
