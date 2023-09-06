package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/DailyAlgorithm"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type DailyAlgorithmRecordSearch struct {
	DailyAlgorithm.DailyAlgorithmRecord
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	UserName       string     `json:"user_name" form:"user_name"`
	request.PageInfo
}
