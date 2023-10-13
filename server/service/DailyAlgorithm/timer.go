package DailyAlgorithm

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/DailyAlgorithm"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"time"
)

// CountDailyAlgorithmRank
// 统计user_name在前30天内的打卡次数，返回一个键值对数组`map[string]int`，每天的打卡最多计算一次
// 关于在定时器中调用，因为是多线程，就要考虑对数据库写入的互斥操作
// 所以下面两个函数必须在对数据库写入时
// 使用读写锁，仅在写入时互斥
//
// 当时只有新开辟的线程会使用这个函数，互斥应该在数据库层面上的互斥
// 这里互斥没有用

func CountDailyAlgorithmRank() (err error) {

	// 定义一个键值对数组，用于存储该月用户的打卡天数
	rank := make(map[string]int)

	// 就算是打卡次数为0，也需要记录一下，为下面rank的map可以更新到所有用户
	// 所以要再遍历一遍所有用户，先加入map中
	var uuidList []string
	uuidList, err = system.GetUserUuidList()
	if err != nil {
		return
	}
	for _, uuid := range uuidList {
		rank[uuid] = 0
	}

	// 获取当前日期
	now := time.Now()

	// 循环遍历过去30天的日期

	for i := 0; i < 30; i++ {
		// 计算当前日期的字符串表示
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		fmt.Println(date)
		// 创建db
		// 在这里是因为循环导包

		var DARs []DailyAlgorithm.DailyAlgorithmRecord
		db := global.GVA_DB.Model(&DailyAlgorithm.DailyAlgorithmRecord{})
		db = db.Where("date = ?", date)

		err = db.Find(&DARs).Error

		//返回错误
		if err != nil {
			return nil
		}

		// 定义一个键值对数组，用于该日用户是否打卡
		// 如果数据库中出现了重复用户和日期打的打卡，有更好的健壮性
		isExist := make(map[string]bool)
		for _, record := range DARs {
			if isExist[record.User_name] == false {
				rank[record.User_name]++
			}
			isExist[record.User_name] = true
		}

	}

	for name, cnt := range rank {
		fmt.Println(name, cnt)
		system.CoverDACount(name, cnt)
	}

	return nil
}

// RemoveCountedOutDate
// 返回哪些用户需要被-1的map
// - 在有人提交记录时，就给该用户今天的打卡次数+1，直接更新数据库中的数据
// - 每天只需要排除第31天之前的所有数据
// - 成功优化了查询效率！！！
// - 但是为了更具有鲁棒性，之前的函数可以保留，在出现错误时重新调用，可以保证排行榜的正确性
// 统计user_name在前30天内的打卡次数，返回一个键值对数组`map[string]int`，每天的打卡最多计算一次
func RemoveCountedOutDate() (err error) {

	// 获取当前日期
	now := time.Now()
	// 计算当前日期的字符串表示
	date := now.AddDate(0, 0, -31).Format("2006-01-02")

	// 创建db
	db := global.GVA_DB.Model(&DailyAlgorithm.DailyAlgorithmRecord{})
	var DARs []DailyAlgorithm.DailyAlgorithmRecord
	//查找
	db = db.Where("date = ?", date)

	//出错返回空map和错误
	err = db.Find(&DARs).Error
	if err != nil {
		return err
	}

	for _, record := range DARs {
		system.UpdateDACount(record.User_name, -1)
	}

	return nil
}

// Timer
// sb代码，在docker中莫名其妙用不了，改用别人写好的函数
func Timer() (err error) {

	//每次运行，调用一次更新30天内的函数，保证数据的可维护性
	err = CountDailyAlgorithmRank()
	if err != nil {
		global.GVA_LOG.Error("DACount Error")
		return err
	}
	global.GVA_LOG.Info("DACount Counted")

	//corn框架中的定时器
	t := timer.NewTimerTask()

	_, err = t.AddTaskByFunc("DAFunc", "@midnight", func() {
		err := RemoveCountedOutDate()
		if err != nil {
			global.GVA_LOG.Error("process counted work occur to error ,recalculate")
			CountDailyAlgorithmRank()
		}
		global.GVA_LOG.Info("update succeed") // 每天打印一遍
	})
	return err
}
