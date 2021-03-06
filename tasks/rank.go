package tasks

import (
	"ADDD_DOUYIN/util"
	"context"
)

// RestartDailyRank 开启新一天的排名
func RestartDailyRank() error {
	ctx := context.Background()
	return util.RedisClient1.Del(ctx, "rank:daily").Err()
}
