package worker

import (
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
)

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
	log.Debug("刷新心跳时间")
	rep, err := repository.NewConfig()
	if err != nil {
		c.errChan <- err
		return
	}
	err = rep.UpdateHeartBeat()
	if err != nil {
		c.errChan <- err
	}
}

//刷新配置
func (c *common) RefreshConfig() {
	//TODO 配置刷新（任务的执行需支持停止，含线程的停止）
	log.Debug("配置刷新（任务的执行需支持停止，含线程的停止）")
}
