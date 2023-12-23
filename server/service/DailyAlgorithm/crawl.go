package DailyAlgorithm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

// Crawl
// 每天爬取一次积分，启动时爬取一次
func Crawl() {
	err := system.GetUserLcRatingList()
	if err != nil {
		global.GVA_LOG.Error(err.Error())
	} else {
		global.GVA_LOG.Info("LcRating Counted")
	}

	_, err = global.GVA_Timer.AddTaskByFunc("DAFunc", "@weekly", func() {
		err := system.GetUserLcRatingList()
		if err != nil {
			global.GVA_LOG.Error("process counted work occur to error ,recalculate")
		}
		global.GVA_LOG.Info("update succeed") // 每天打印一遍
	})

}
