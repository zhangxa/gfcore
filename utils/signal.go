package utils

import (
	"sync"
)

// Signal 信号类
var Signal = new(signal)

type signal struct{}

// StopSignal 停止信号值
type StopSignal int32

type exitWait struct {
	mutex          sync.Mutex
	wg             *sync.WaitGroup
	deferFunctions []func()
	stopSignList   []chan StopSignal
}

var exitWaitHandler *exitWait

func init() {
	exitWaitHandler = &exitWait{
		wg: &sync.WaitGroup{},
	}
}

// ExitWaitFunDo 退出后等待处理完成
//
//	@Title  退出后等待处理完成
//	@Description
//	@Author  Ms <133814250@qq.com>
//	@Param   doFun
func (util *signal) ExitWaitFunDo(doFun func()) {
	exitWaitHandler.wg.Add(1)
	defer exitWaitHandler.wg.Done()
	if doFun != nil {
		doFun()
	}
}

// AppDefer 应用退出后置操作
//
//	@Title  应用退出后置操作
//	@Description
//	@Author  Ms <133814250@qq.com>
//	@Param   deferFun
func (util *signal) AppDefer(deferFun ...func()) {
	exitWaitHandler.mutex.Lock()
	defer exitWaitHandler.mutex.Unlock()
	for _, funcItem := range deferFun {
		if funcItem != nil {
			exitWaitHandler.deferFunctions = append(exitWaitHandler.deferFunctions, funcItem)
		}
	}
}

// ListenStop 订阅app退出信号
//
//	@Title  订阅app退出信号
//	@Description
//	@Author  Ms <133814250@qq.com>
//	@Param   stopSig
func (util *signal) ListenStop(stopSig chan StopSignal) {
	exitWaitHandler.mutex.Lock()
	defer exitWaitHandler.mutex.Unlock()

	exitWaitHandler.stopSignList = append(exitWaitHandler.stopSignList, stopSig)
}
