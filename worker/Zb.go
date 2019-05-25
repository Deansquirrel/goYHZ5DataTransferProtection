package worker

import (
	"errors"
	"fmt"
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
	log.Debug("恢复门店设置")
	repZb, err := repository.NewRepZb()
	if err != nil {
		w.errChan <- err
		return
	}
	rep, err := repository.NewConfig()
	if err != nil {
		w.errChan <- err
	}
	for {
		lstOpr, err := repZb.GetLstMdSetOpr()
		if err != nil {
			w.errChan <- err
			return
		}
		if lstOpr == nil {
			return
		}
		if lstOpr.OprType != 1 && lstOpr.OprType != 2 && lstOpr.OprType != 3 {
			errMsg := fmt.Sprintf("opr type err,exp 1 or 2 or 3,got %d", lstOpr.OprType)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
			return
		}
		err = rep.UpdateMdSet(lstOpr)
		if err != nil {
			w.errChan <- err
			return
		}
		err = repZb.DelMdSetOpr(lstOpr.OprSn)
		if err != nil {
			w.errChan <- err
			return
		}
	}
}

//恢复财务公司设置
func (w *zb) RestoreCwGsSet() {
	log.Debug("恢复财务公司设置")
	repZb, err := repository.NewRepZb()
	if err != nil {
		w.errChan <- err
		return
	}
	rep, err := repository.NewConfig()
	if err != nil {
		w.errChan <- err
	}
	for {
		lstOpr, err := repZb.GetLstCwGsSetOpr()
		if err != nil {
			w.errChan <- err
			return
		}
		if lstOpr == nil {
			return
		}
		if lstOpr.OprType != 1 && lstOpr.OprType != 2 && lstOpr.OprType != 3 {
			errMsg := fmt.Sprintf("opr type err,exp 1 or 2 or 3,got %d", lstOpr.OprType)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
			return
		}
		err = rep.UpdateCwGsSet(lstOpr)
		if err != nil {
			w.errChan <- err
			return
		}
		err = repZb.DelCwGsSetOpr(lstOpr.OprSn)
		if err != nil {
			w.errChan <- err
			return
		}
	}
}

//恢复门店隶属财务公司关系设置
func (w *zb) RestoreMdCwGsRef() {
	log.Debug("恢复门店隶属财务公司关系设置")
	repZb, err := repository.NewRepZb()
	if err != nil {
		w.errChan <- err
		return
	}
	rep, err := repository.NewConfig()
	if err != nil {
		w.errChan <- err
	}
	for {
		lstOpr, err := repZb.GetLstMdCwGsRefOpr()
		if err != nil {
			w.errChan <- err
			return
		}
		if lstOpr == nil {
			return
		}
		if lstOpr.OprType != 1 && lstOpr.OprType != 2 && lstOpr.OprType != 3 {
			errMsg := fmt.Sprintf("opr type err,exp 1 or 2 or 3,got %d", lstOpr.OprType)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
			return
		}
		err = rep.UpdateMdCwGsRef(lstOpr)
		if err != nil {
			w.errChan <- err
			return
		}
		err = repZb.DelMdCwGsRefOpr(lstOpr.OprSn)
		if err != nil {
			w.errChan <- err
			return
		}
	}
}
