package util

import "github.com/liangyaopei/bloom"

//维护一个全局bloomFilter，仅在服务器开启时在内存中
//因此，当服务器重启后，原bloomfilter里的数据就丢失了,重启后的首位注册用户才会去查数据库
var Filter *bloom.Filter

func InitBloomFilter() {
	filter := bloom.New(1024, 3, false)
	Filter = filter
}
