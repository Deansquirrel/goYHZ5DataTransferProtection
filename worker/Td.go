package worker

import (
	"fmt"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
	"strconv"
)

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
			}
		}()
	}
}

//恢复门店营业日开闭店记录传递时间
func (w *td) RestoreMdYyStateTransTime() {
	// TODO 恢复门店营业日开闭店记录传递时间
	log.Debug("恢复门店营业日开闭店记录传递时间")
}
