package worker

import (
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
)

type zb struct {
	errChan chan<- error
}

func NewWorkerZb(errChan chan<- error) *zb {
	return &zb{
		errChan: errChan,
	}
}

//刷新待恢复信道数据总量
func (w *zb) RefreshWaitRestoreDataCount() {
	log.Debug("刷新待恢复信道数据总量")
	repZb, err := repository.NewRepZb()
	if err != nil {
		w.errChan <- err
		return
	}
	repConfig, err := repository.NewConfig()
	if err != nil {
		w.errChan <- err
		return
	}
	waitRestoreDataCount, err := repZb.GetWaitRestoreDataCount()
	if err != nil {
		w.errChan <- err
		return
	}
	err = repConfig.AddWaitRestoreDataCountRecord(waitRestoreDataCount)
	if err != nil {
		w.errChan <- err
		return
	}
	for i := 1; i <= 4; i++ {
		ri := i
		go func() {
			count, uTime, err := repZb.GetChanWaitRestoreDataCount(ri)
			if err != nil {
				w.errChan <- err
			} else {
				err = repConfig.UpdateWaitRestoreDataCount(waitRestoreDataCount.BatchNo, ri, uTime, count)
				if err != nil {
					w.errChan <- err
				}
			}
		}()
	}
}

//恢复门店营业日开闭店记录恢复时间
func (w *zb) RestoreMdYyStateRestoreTime() {
	// TODO 恢复门店营业日开闭店记录恢复时间
	log.Debug("恢复门店营业日开闭店记录恢复时间")
}

//恢复门店营业日开闭店记录
func (w *zb) RestoreMdYyState() {
	// TODO 恢复门店营业日开闭店记录
	log.Debug("恢复门店营业日开闭店记录")
}

//恢复门店设置
func (w *zb) RestoreMdSet() {
	//TODO 恢复门店设置
	log.Debug("恢复门店设置")
}

//恢复财务公司设置
func (w *zb) RestoreCwGsSet() {
	//TODO 恢复财务公司设置
	log.Debug("恢复财务公司设置")
}

//恢复门店隶属财务公司关系设置
func (w *zb) RestoreMdCwGsRef() {
	//TODO 恢复门店隶属财务公司关系设置
	log.Debug("恢复门店隶属财务公司关系设置")
}
