package worker

import log "github.com/Deansquirrel/goToolLog"

type common struct {
	errChan chan<- error
}

func NewCommon(errChan chan<- error) *common {
	return &common{
		errChan: errChan,
	}
}

//刷新心跳时间
func (c *common) RefreshHeartBeat() {
	//TODO 刷新心跳时间
	log.Debug("刷新心跳时间")
}
