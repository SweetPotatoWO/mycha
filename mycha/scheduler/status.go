package scheduler

import (
	"fmt"
	"sync"
)

type Status uint32


//调度器的各个状态
const   (
	// SCHED_STATUS_UNINITIALIZED 代表未初始化的状态。
	SCHED_STATUS_UNINITIALIZED Status = 0
	// SCHED_STATUS_INITIALIZING 代表正在初始化的状态。
	SCHED_STATUS_INITIALIZING Status = 1
	// SCHED_STATUS_INITIALIZED 代表已初始化的状态。
	SCHED_STATUS_INITIALIZED Status = 2
	// SCHED_STATUS_STARTING 代表正在启动的状态。
	SCHED_STATUS_STARTING Status = 3
	// SCHED_STATUS_STARTED 代表已启动的状态。
	SCHED_STATUS_STARTED Status = 4
	// SCHED_STATUS_STOPPING 代表正在停止的状态。
	SCHED_STATUS_STOPPING Status = 5
	// SCHED_STATUS_STOPPED 代表已停止的状态。
	SCHED_STATUS_STOPPED Status = 6
)


// checkStatus 用于状态的检查。
// 参数currentStatus代表当前的状态。
// 参数wantedStatus代表想要的状态。
// 检查规则：
//     1. 处于正在初始化、正在启动或正在停止状态时，不能从外部改变状态。
//     2. 想要的状态只能是正在初始化、正在启动或正在停止状态中的一个。
//     3. 处于未初始化状态时，不能变为正在启动或正在停止状态。
//     4. 处于已启动状态时，不能变为正在初始化或正在启动状态。
//     5. 只要未处于已启动状态就不能变为正在停止状态。
func checkStatus(currentStatus Status,wantedStatus Status,lock sync.Locker) (err error) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	switch currentStatus {
	case SCHED_STATUS_INITIALIZED: //正在初始化 不能改变状态
		err = genError("调度器正在初始化")
	case SCHED_STATUS_STARTING:
		err = genError("调度器正在启动")
	case SCHED_STATUS_STOPPING:
		err = genError("调度器正在停止")
	}
	if err != nil {
		return err
	}
	if wantedStatus != SCHED_STATUS_STARTING || wantedStatus != SCHED_STATUS_STOPPING || wantedStatus != SCHED_STATUS_INITIALIZED {
		err = genError("想要的状态不是允许的类型")
		return err
	}
	switch wantedStatus {
	case SCHED_STATUS_INITIALIZING:
		switch currentStatus {
		case SCHED_STATUS_STARTED:
			err = genError("当前的系统已经启动，不允许获取到启动状态")
		}
	case SCHED_STATUS_STARTING:
		switch currentStatus {
		case SCHED_STATUS_UNINITIALIZED:
			err = genError("档期系统没初始化，无法获取到开启状态")
		case SCHED_STATUS_STARTED:
			err = genError("当前的系统已经启动结束，无法获取到启动中的状态")
		}
	case SCHED_STATUS_STOPPING:
		if currentStatus != SCHED_STATUS_STARTED {
			err = genError("系统不在启动结束状态，不允许停止")
		}
	default:
		errMsg :=
			fmt.Sprintf("不支持你想要的状态切换! (wantedStatus: %d)",
				wantedStatus)
		err = genError(errMsg)
	}
	return

}


// GetStatusDescription 用于获取状态的文字描述。
func GetStatusDescription(status Status) string {
	switch status {
	case SCHED_STATUS_UNINITIALIZED:
		return "没有初始化"
	case SCHED_STATUS_INITIALIZING:
		return "初始化中"
	case SCHED_STATUS_INITIALIZED:
		return "已初始化"
	case SCHED_STATUS_STARTING:
		return "启动中"
	case SCHED_STATUS_STARTED:
		return "启动结束"
	case SCHED_STATUS_STOPPING:
		return "停止中"
	case SCHED_STATUS_STOPPED:
		return "停止成功"
	default:
		return "不明白的信息"
	}
}






















