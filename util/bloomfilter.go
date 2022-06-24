package util

import "github.com/liangyaopei/bloom"

//维护一个全局bloomFilter，仅在服务器开启时存在于内存中
var Filter *bloom.Filter

func InitBloomFilter() {
	filter := bloom.New(1024, 3, false)
	Filter = filter
}
