package worker

import log "github.com/Deansquirrel/goToolLog"

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
}

//恢复门店营业日开闭店记录传递时间
func (w *td) RestoreMdYyStateTransTime() {
	// TODO 恢复门店营业日开闭店记录传递时间
	log.Debug("恢复门店营业日开闭店记录传递时间")
}
