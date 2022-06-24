package util

import "github.com/charizer/snowflake"

var worker *snowflake.Worker

func init() {
	worker, _ = snowflake.NewWorker(1)
}

func NextId() int64 {
	return worker.Next()
}
