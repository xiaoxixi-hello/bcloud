package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

func InitPanID() string {
	parse, _ := time.Parse("2006-01-02", "2020-02-02")

	snowflake.Epoch = parse.UnixNano() / 1e6
	node, _ := snowflake.NewNode(1000)
	return time.Now().Format("20060102") + node.Generate().String()
}
