package DailyAlgorithm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/DailyAlgorithm"
	DailyAlgorithmReq "github.com/flipped-aurora/gin-vue-admin/server/model/DailyAlgorithm/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"time"
)

type DailyAlgorithmRecordService struct {
}

// CreateDailyAlgorithmRecord 创建DailyAlgorithmRecord记录
// Author [piexlmax](https://github.com/piexlmax)
func (DARService *DailyAlgorithmRecordService) CreateDailyAlgorithmRecord(DAR *DailyAlgorithm.DailyAlgorithmRecord) (err error) {
	tx := global.GVA_DB.Begin() // 开启事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 发生错误时回滚事务
		}
	}()
	// 操作
	err = global.GVA_DB.Create(DAR).Error

	if err != nil {
		tx.Rollback() // 操作失败，回滚事务
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
	}
	err = tx.Commit().Error // 提交事务

	return err
}

// DeleteDailyAlgorithmRecord 删除DailyAlgorithmRecord记录
// Author [piexlmax](https://github.com/piexlmax)
func (DARService *DailyAlgorithmRecordService) DeleteDailyAlgorithmRecord(DAR DailyAlgorithm.DailyAlgorithmRecord) (err error) {
	err = global.GVA_DB.Delete(&DAR).Error
	return err
}

// DeleteDailyAlgorithmRecordByIds 批量删除DailyAlgorithmRecord记录
// Author [piexlmax](https://github.com/piexlmax)
func (DARService *DailyAlgorithmRecordService) DeleteDailyAlgorithmRecordByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]DailyAlgorithm.DailyAlgorithmRecord{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateDailyAlgorithmRecord 更新DailyAlgorithmRecord记录
// Author [piexlmax](https://github.com/piexlmax)
func (DARService *DailyAlgorithmRecordService) UpdateDailyAlgorithmRecord(DAR DailyAlgorithm.DailyAlgorithmRecord) (err error) {
	err = global.GVA_DB.Save(&DAR).Error
	return err
}

// GetDailyAlgorithmRecord 根据id获取DailyAlgorithmRecord记录
// Author [piexlmax](https://github.com/piexlmax)
func (DARService *DailyAlgorithmRecordService) GetDailyAlgorithmRecord(id uint) (DAR DailyAlgorithm.DailyAlgorithmRecord, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&DAR).Error
	return
}

// GetCoverDailyAlgorithmRecord 根据date获取DailyAlgorithmRecord记录
// Author [CFDDFC](https://github.com/cfddd)
func (DARService *DailyAlgorithmRecordService) GetCoverDailyAlgorithmRecord(date time.Time, user_name string) (DAR DailyAlgorithm.DailyAlgorithmRecord, err error) {
	dateString := date.Format("2006-01-02")
	//fmt.Println(dateString)
	//fmt.Println(user_name)
	err = global.GVA_DB.Where("date = ?", dateString).Where("user_name = ?", user_name).First(&DAR).Error
	return
}

// GetDailyAlgorithmRecordInfoList 分页获取DailyAlgorithmRecord记录
// Author [piexlmax](https://github.com/piexlmax)
// 两种用法：
// 1.获取在每日算法打卡记录列表里面使用，info的User_name为空，查询所有的打卡记录
// 2.获取某一个用户的全部打卡记录，info的User_name不为空，在查询时多查询一次，先找到uuid再在记录里面查询
func (DARService *DailyAlgorithmRecordService) GetDailyAlgorithmRecordInfoList(info DailyAlgorithmReq.DailyAlgorithmRecordSearch) (list []DailyAlgorithm.DailyAlgorithmRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var uuid string

	if info.User_name != "" {
		// 创建db
		var user system.SysUser
		db := global.GVA_DB.Model(&system.SysUser{})
		err = db.Where("nick_name = ?", info.User_name).Find(&user).Error
		if err != nil {
			return
		}
		uuid = user.UUID.String()
	}

	// 创建db
	db := global.GVA_DB.Model(&DailyAlgorithm.DailyAlgorithmRecord{})
	var DARs []DailyAlgorithm.DailyAlgorithmRecord

	// Add order by id desc
	db = db.Order("id DESC")

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Date != nil {
		db = db.Where("date = ?", info.Date)
	}
	if uuid != "" {
		db = db.Where("user_name = ?", uuid)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&DARs).Error
	return DARs, total, err
}
