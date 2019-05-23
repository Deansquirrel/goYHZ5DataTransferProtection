package worker

import (
	"encoding/json"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
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
	//TODO 刷新门店数据信道传递状态
	log.Debug("刷新门店数据信道传递状态")
	rep, err := repository.NewRepTd()
	if err != nil {
		w.errChan <- err
		return
	}
	sList, err := rep.GetMdDataTransState()
	if err != nil {
		w.errChan <- err
	}

	if sList != nil {
		for _, s := range sList {
			dataStr, err := json.Marshal(s)
			if err != nil {
				log.Debug(err.Error())
			} else {
				log.Debug(string(dataStr))
			}
		}
	}
	//TODO Save sList
}

//恢复门店营业日开闭店记录传递时间
func (w *td) RestoreMdYyStateTransTime() {
	// TODO 恢复门店营业日开闭店记录传递时间
	log.Debug("恢复门店营业日开闭店记录传递时间")
}
