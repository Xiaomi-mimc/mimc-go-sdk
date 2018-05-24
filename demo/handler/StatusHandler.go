package handler

type StatusHandler struct {
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

func (this StatusHandler) HandleChange(isOnline bool, errType, errReason, errDescription *string) {
	if isOnline {
		logger.Info("status changed: online.")
	} else {
		logger.Info("status changed: offline.")
	}
}
