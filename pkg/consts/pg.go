package consts

import "time"

const (
	MaxIdleConn = 3
	MaxOpenConn = 90
	MaxLifetime = 1 * time.Hour
	MaxIdleTime = 2 * time.Second
)
