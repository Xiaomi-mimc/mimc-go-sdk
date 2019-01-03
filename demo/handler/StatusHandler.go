package handler

type StatusHandler struct {
	appAccount string
}

func NewStatusHandler(appAccount string) *StatusHandler {
	return &StatusHandler{appAccount}
}

func (this StatusHandler) HandleChange(isOnline bool, errType, errReason, errDescription *string) {
	if isOnline {
		logger.Info("[%v] status changed: online.", this.appAccount)
	} else {
		logger.Info("[%v] status changed: offline，errType:%v, errReason:%v, errDes:%v", this.appAccount, *errType, *errReason, *errDescription)
	}
}
