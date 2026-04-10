package utils

import (
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// InitSnowflake 初始化雪花算法
func InitSnowflake() {
	once.Do(func() {
		// 节点编号 0-1023，单机默认1即可
		n, err := snowflake.NewNode(1)
		if err != nil {
			panic("雪花算法初始化失败: " + err.Error())
		}
		node = n
	})
}

// GenStringID 生成 String 类型的唯一ID
func GenStringID() string {
	return strconv.FormatInt(node.Generate().Int64(), 10)
}
