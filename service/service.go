package service

import (
	"context"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/worker"
)

//启动服务内容
func StartService() error {
	log.Debug("StartService")

	ch := make(chan error)
	commWorker := worker.NewCommon(ch)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case err := <-ch:
				if err != nil {
					log.Debug(err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	//配置刷新
	commWorker.RefreshConfig()

	////Test 显示任务运行情况
	//{
	//	go func() {
	//		c := cron.New()
	//		err := c.AddFunc("0 * * * * ?", showTaskState)
	//		if err != nil {
	//			log.Debug(err.Error())
	//		} else {
	//			c.Start()
	//		}
	//	}()
	//}
	cancel()
	return nil
}

////显示任务状态
//func showTaskState() {
//	idList := global.TaskList.GetIdList()
//	sort.Strings(idList)
//
//	for _, k := range idList {
//		testPrintS(global.TaskList.GetObject(k).(*object.TaskState))
//	}
//}

////输出TaskState
//func testPrintS(ts *object.TaskState) {
//	if ts == nil {
//		log.Debug("ts is nil")
//		return
//	}
//	str, err := json.Marshal(ts)
//	if err != nil {
//		log.Debug(fmt.Sprintf("%s - %s", ts.Key, err.Error()))
//	} else {
//		log.Debug(fmt.Sprintf("%s - %s", ts.Key, str))
//	}
//	c := ts.Cron
//	if c != nil {
//		if len(c.Entries()) > 0 {
//			e := c.Entries()[0]
//			log.Debug(fmt.Sprintf("Prev time: %s,Next time: %s",
//				goToolCommon.GetDateTimeStr(e.Prev),
//				goToolCommon.GetDateTimeStr(e.Next)))
//		}
//	}
//}
