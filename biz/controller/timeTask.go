package controller

import (
	"TikTok/biz/service/timeTaskService"
)

func TimeTaskExec() {
	timeTaskService.SyncFavorDb()
}
