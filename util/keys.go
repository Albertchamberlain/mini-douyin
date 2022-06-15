package util

import (
	"fmt"
	"strconv"
)

const (
	// DailyRankKey 每日排行
	DailyRankKey = "rank:daily"
)

// VideoViewKey 视频点击数的key
// view:video:1 -> 100
// view:video:2 -> 150
func VideoViewKey(id int64) string {
	return fmt.Sprintf("view:video:%s", strconv.Itoa(int(id)))
}
