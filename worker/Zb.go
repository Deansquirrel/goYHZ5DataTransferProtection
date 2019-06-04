package worker

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
	"sync"
)

var syncLockRefreshWaitRestoreDataCount sync.Mutex
var syncLockRestoreMdYyStateRestoreTime sync.Mutex
var syncLockRestoreMdYyState sync.Mutex
var syncLockRestoreMdSet sync.Mutex
var syncLockRestoreCwGsSet sync.Mutex
var syncLockRestoreMdCwGsRef sync.Mutex

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

	key := object.TaskKeyRefreshWaitRestoreDataCount
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRefreshWaitRestoreDataCount.Lock()
	defer func() {
		syncLockRefreshWaitRestoreDataCount.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

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
	log.Debug("恢复门店营业日开闭店记录恢复时间")

	key := object.TaskKeyRestoreMdYyStateRestoreTime
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRestoreMdYyStateRestoreTime.Lock()
	defer func() {
		syncLockRestoreMdYyStateRestoreTime.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

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
		lstOpr, err := repZb.GetLstMdYyStateRestoreTimeOpr()
		if err != nil {
			w.errChan <- err
			return
		}
		if lstOpr == nil {
			w.errChan <- nil
			return
		}
		if lstOpr.OprType != 1 && lstOpr.OprType != 2 {
			errMsg := fmt.Sprintf("opr type err,exp 1 or 2 ,got %d", lstOpr.OprType)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
			return
		}
		err = rep.UpdateMdYyStateRestoreTime(lstOpr)
		if err != nil {
			w.errChan <- err
			return
		}
		err = repZb.DelMdYyStateRestoreTimeOpr(lstOpr.OprSn)
		if err != nil {
			w.errChan <- err
			return
		}
	}
}

//恢复门店营业日开闭店记录
func (w *zb) RestoreMdYyState() {
	log.Debug("恢复门店营业日开闭店记录")

	key := object.TaskKeyRestoreMdYyState
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRestoreMdYyState.Lock()
	defer func() {
		syncLockRestoreMdYyState.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

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
		lstOpr, err := repZb.GetLstMdYyStateOpr()
		if err != nil {
			w.errChan <- err
			return
		}
		if lstOpr == nil {
			w.errChan <- nil
			return
		}
		if lstOpr.OprType != 1 && lstOpr.OprType != 2 {
			errMsg := fmt.Sprintf("opr type err,exp 1 or 2 ,got %d", lstOpr.OprType)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
			return
		}
		err = rep.UpdateMdYyState(lstOpr)
		if err != nil {
			w.errChan <- err
			return
		}
		err = repZb.DelMdYyStateOpr(lstOpr.OprSn)
		if err != nil {
			w.errChan <- err
			return
		}
	}
}

//恢复门店设置
func (w *zb) RestoreMdSet() {
	log.Debug("恢复门店设置")

	key := object.TaskKeyRestoreMdSet
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRestoreMdSet.Lock()
	defer func() {
		syncLockRestoreMdSet.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

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
			w.errChan <- nil
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

	key := object.TaskKeyRestoreCwGsSet
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRestoreCwGsSet.Lock()
	defer func() {
		syncLockRestoreCwGsSet.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

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
			w.errChan <- nil
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

	key := object.TaskKeyRestoreMdCwGsRef
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRestoreMdCwGsRef.Lock()
	defer func() {
		syncLockRestoreMdCwGsRef.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

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
			w.errChan <- nil
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
