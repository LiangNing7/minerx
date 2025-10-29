package rid

import (
	"github.com/LiangNing7/goutils/pkg/id"
)

const defaultABC = "abcdefghijklmnopqrstuvwxyz1234567890"

type ResourceID string

const (
	// ID for the user resource in minerx-usercenter.
	User ResourceID = "user"
	// ID for the order resource in minerx-fakeserver.
	Order ResourceID = "order"
	// ID for the cronjob resource in minerx-nightwatch.
	CronJob ResourceID = "cronjob"
	// ID for the job resource in minerx-nightwatch.
	Job ResourceID = "job"
)

// String 将资源标识符转换为字符串.
func (rid ResourceID) String() string {
	return string(rid)
}

// New 创建带前缀的唯一标识符.
func (rid ResourceID) New(counter uint64) string {
	// 使用自定义选项生成唯一标识符
	uniqueStr := id.NewCode(
		counter,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return rid.String() + "-" + uniqueStr
}
