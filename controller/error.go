package controller

type ControllerError struct {
	F   string
	Msg string
	Err error
}

func (e ControllerError) Error() string {
	return "Controller: " + e.F + " " + e.Msg + "\n" + e.Err.Error()
}
