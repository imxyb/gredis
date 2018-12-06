package consts

import "errors"

var (
	ErrCompareType = errors.New("请勿对比非字符串对象")
)
