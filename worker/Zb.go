package worker

import log "github.com/Deansquirrel/goToolLog"

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
	//TODO 刷新待恢复信道数据总量
	log.Debug("刷新待恢复信道数据总量")
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
