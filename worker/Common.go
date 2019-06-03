package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/global"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
	"github.com/robfig/cron"
	"strings"
	"sync"
)

const (
	//TODO 钉钉机器人消息Key待参数化
	webHook = "3dd601535aa95027b92e78b8a820ba62be5069293092a302d8c17ef63e095cac"
)

var syncLock sync.Mutex

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
		return
	}
	c.errChan <- nil
}

//刷新配置
func (c *common) RefreshConfig() {
	log.Debug("配置刷新")
	rep, err := repository.NewConfig()
	if err != nil {
		c.errChan <- err
		return
	}
	idList := global.TaskKeyList
	for _, id := range idList {
		t := global.TaskList.GetObject(string(id))
		if t == nil {
			c.startService(id, c.ErrHandle)
			continue
		}
		ts := t.(*object.TaskState)
		newTs, err := rep.GetTaskCronByKey(ts.Key)
		if err != nil {
			c.errChan <- err
			return
		}
		c.errChan <- nil
		if newTs.Cron != ts.CronStr {
			//停止工作线程
			c.stopWorker(id)
			//启动工作线程
			c.startService(id, c.ErrHandle)
		}
	}
}

//根据Key，启动制定服务
func (c *common) startService(key object.TaskKey, errHandle func(err error)) {
	switch key {
	case object.TaskKeyHeartBeat:
		c.startHeartBeat(errHandle)
	case object.TaskKeyRefreshMdDataTransState:
		c.startRefreshMdDataTransState(errHandle)
	case object.TaskKeyRestoreMdYyStateTransTime:
		c.startRestoreMdYyStateTransTime(errHandle)
	case object.TaskKeyRefreshWaitRestoreDataCount:
		c.startRefreshWaitRestoreDataCount(errHandle)
	case object.TaskKeyRestoreMdYyStateRestoreTime:
		c.startRestoreMdYyStateRestoreTime(errHandle)
	case object.TaskKeyRestoreMdYyState:
		c.startRestoreMdYyState(errHandle)
	case object.TaskKeyRestoreMdSet:
		c.startRestoreMdSet(errHandle)
	case object.TaskKeyRestoreCwGsSet:
		c.startRestoreCwGsSet(errHandle)
	case object.TaskKeyRestoreMdCwGsRef:
		c.startRestoreMdCwGsRef(errHandle)
	case object.TaskKeyRefreshConfig:
		c.StartRefreshConfig(errHandle)
	default:
		errMsg := fmt.Sprintf("unknow task key: %s", key)
		log.Error(errMsg)
		c.errChan <- errors.New(errMsg)
	}
}

//刷新配置
func (c *common) StartRefreshConfig(errHandle func(err error)) {
	ch := make(chan error)
	comm := NewCommon(ch)
	c.startWorker(object.TaskKeyRefreshConfig, comm.RefreshConfig, ch, errHandle)
}

//Common（心跳）
func (c *common) startHeartBeat(errHandle func(err error)) {
	ch := make(chan error)
	comm := NewCommon(ch)
	c.startWorker(object.TaskKeyHeartBeat, comm.RefreshHeartBeat, ch, errHandle)
}

//通道
func (c *common) startRefreshMdDataTransState(errHandle func(err error)) {
	ch := make(chan error)
	workerTd := NewWorkerTd(ch)
	c.startWorker(object.TaskKeyRefreshMdDataTransState, workerTd.RefreshMdDataTransState, ch, errHandle)
}

//通道
func (c *common) startRestoreMdYyStateTransTime(errHandle func(err error)) {
	ch := make(chan error)
	workerTd := NewWorkerTd(ch)
	c.startWorker(object.TaskKeyRestoreMdYyStateTransTime, workerTd.RestoreMdYyStateTransTime, ch, errHandle)
}

//总部库
func (c *common) startRefreshWaitRestoreDataCount(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRefreshWaitRestoreDataCount, w.RefreshWaitRestoreDataCount, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdYyStateRestoreTime(errHandle func(error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdYyStateRestoreTime, w.RestoreMdYyStateRestoreTime, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdYyState(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdYyState, w.RestoreMdYyState, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdSet(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdSet, w.RestoreMdSet, ch, errHandle)
}

//总部库
func (c *common) startRestoreCwGsSet(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreCwGsSet, w.RestoreCwGsSet, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdCwGsRef(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdCwGsRef, w.RestoreMdCwGsRef, ch, errHandle)
}

//启动工作线程
func (c *common) startWorker(key object.TaskKey, cmd func(), ch chan error, errHandle func(err error)) {
	s := &object.TaskState{
		Key:     key,
		Cron:    nil,
		CronStr: "",
		Running: false,
		Err:     nil,
	}
	s.Ctx, s.Cancel = context.WithCancel(context.Background())
	{
		syncLock.Lock()
		task := global.TaskList.GetObject(string(s.Key))
		if task != nil {
			c.stopWorker(s.Key)
		}
		global.TaskList.Register() <- goToolCommon.NewObject(string(s.Key), s)
		syncLock.Unlock()
	}

	go func() {
		defer func() {
			close(ch)
		}()
		for {
			select {
			case err := <-ch:
				if err != nil {
					errHandle(err)
				}
				if s.Err != nil && err == nil {
					c.sendMsg(fmt.Sprintf("Task %s is resume", key))
				}
				s.Err = err
			case <-s.Ctx.Done():
				//c.stopWorker(s.Key)
				return
			case <-global.Ctx.Done():
				//c.stopWorker(s.Key)
				return
			}
		}
	}()

	cronStr, err := c.getTaskCron(s.Key)
	if err != nil {
		s.CronStr = ""
		c.errChan <- err
		return
	}
	cronStr = strings.Trim(cronStr, " ")
	if cronStr == "" {
		c.errChan <- errors.New("cron is empty")
		return
	}
	s.CronStr = cronStr
	cr := cron.New()
	err = cr.AddFunc(s.CronStr, cmd)
	if err != nil {
		c.errChan <- err
		return
	}

	log.Debug(fmt.Sprintf("start worker %s cron %s", key, cronStr))

	cr.Start()
	s.Running = true
	s.Cron = cr
}

//停止工作线程
func (c *common) stopWorker(key object.TaskKey) {
	t := global.TaskList.GetObject(string(key))
	if t == nil {
		return
	}
	global.TaskList.Unregister() <- string(key)
	log.Debug(fmt.Sprintf("stop worker %s", key))
	ts := t.(*object.TaskState)
	if ts.Running {
		ts.Running = false
		ts.Cron.Stop()
		ts.Cancel()
	}
}

//根据key获取任务执行Cron时间公式
func (c *common) getTaskCron(key object.TaskKey) (string, error) {
	repConfig, err := repository.NewConfig()
	if err != nil {
		return "", err
	}
	taskCron, err := repConfig.GetTaskCronByKey(key)
	if err != nil {
		return "", err
	}
	if taskCron == nil {
		return "", errors.New("get taskCron err: return is nil")
	}
	return taskCron.Cron, nil
}

func (c *common) ErrHandle(err error) {
	if err == nil {
		return
	}
	c.sendMsg(err.Error())
}

func (c *common) sendMsg(msg string) {
	dt := object.NewDingTalkRobot(&object.DingTalkRobotConfigData{
		FWebHookKey: webHook,
		FAtMobiles:  "",
		FIsAtAll:    0,
	})
	sendErr := dt.SendMsg(msg)
	if sendErr != nil {
		log.Error(fmt.Sprintf("send msg,msg: %s, error: %s", msg, sendErr.Error()))
	}
}
