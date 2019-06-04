package worker

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
	"strconv"
	"sync"
)

var syncLockRefreshMdDataTransState sync.Mutex
var syncLockRestoreMdYyStateTransTime sync.Mutex

type td struct {
	errChan chan<- error
}

func NewWorkerTd(errChan chan<- error) *td {
	return &td{
		errChan: errChan,
	}
}

//刷新门店数据信道传递状态
func (w *td) RefreshMdDataTransState() {
	log.Debug("刷新门店数据信道传递状态")

	key := object.TaskKeyRefreshMdDataTransState
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRefreshMdDataTransState.Lock()
	defer func() {
		syncLockRefreshMdDataTransState.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

	repTd, err := repository.NewRepTd()
	if err != nil {
		w.errChan <- err
		return
	}
	rep, err := repository.NewConfig()
	if err != nil {
		w.errChan <- err
		return
	}
	sList, err := repTd.GetMdDataTransState()
	if err != nil {
		w.errChan <- err
		return
	}
	err = rep.AddMdDataTransState(sList)
	if err != nil {
		w.errChan <- err
		return
	}
	for _, d := range sList {
		cd := d
		go func() {
			tableName := fmt.Sprintf("zl" + object.TdUserCode + cd.MdCode + "tfupv3_" + strconv.Itoa(cd.ChanId))
			rowCount, err := repTd.GetDataRowCount(tableName)
			if err != nil {
				w.errChan <- err
				return
			}
			err = rep.UpdateMdDataTransState(cd.BatchNo, cd.MdCode, cd.ChanId, rowCount)
			if err != nil {
				w.errChan <- err
				return
			}
		}()
	}
	w.errChan <- nil
}

//恢复门店营业日开闭店记录传递时间
func (w *td) RestoreMdYyStateTransTime() {
	log.Debug("恢复门店营业日开闭店记录传递时间")

	key := object.TaskKeyRestoreMdYyStateTransTime
	guid := goToolCommon.Guid()
	log.Debug(fmt.Sprintf("task %s[%s] start", key, guid))
	syncLockRestoreMdYyStateTransTime.Lock()
	defer func() {
		syncLockRestoreMdYyStateTransTime.Unlock()
		//错误处理（panic）
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("task recover get err: %s", err)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
		}
		log.Debug(fmt.Sprintf("task %s[%s] complete", key, guid))
	}()

	repTd, err := repository.NewRepTd()
	if err != nil {
		w.errChan <- err
		return
	}
	rep, err := repository.NewConfig()
	if err != nil {
		w.errChan <- err
	}
	for {
		lstOpr, err := repTd.GetLstMdYyStateTransTimeTdOpr()
		if err != nil {
			w.errChan <- err
			return
		}
		if lstOpr == nil {
			w.errChan <- nil
			return
		}
		if lstOpr.OprType != 1 && lstOpr.OprType != 2 && lstOpr.OprType != 3 && lstOpr.OprType != 4 {
			errMsg := fmt.Sprintf("opr type err,exp 1 or 2 or 3 or 4,got %d", lstOpr.OprType)
			log.Error(errMsg)
			w.errChan <- errors.New(errMsg)
			return
		}
		err = rep.UpdateMdYyStateTransTimeTd(lstOpr)
		if err != nil {
			w.errChan <- err
			return
		}
		err = repTd.DelMdYyStateTransTimeTdOpr(lstOpr.OprSn)
		if err != nil {
			w.errChan <- err
			return
		}
	}
}
