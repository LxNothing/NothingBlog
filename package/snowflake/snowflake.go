package snowflake

// 雪花算法 - 使用推特的snowflake包

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var sf_node *sf.Node

// startTime - 雪花算法的起始时间
func Init(startTime string, machineId int64) (err error) {

	var start_time time.Time

	start_time, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	sf.Epoch = start_time.UnixNano() / 1000000 // 得到毫秒值

	sf_node, err = sf.NewNode(machineId)
	return
}

func GetNextId() sf.ID {
	return sf_node.Generate()
}
