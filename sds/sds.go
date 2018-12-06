package sds

import (
	"unicode/utf8"
	"errors"
	"math"
	"strings"
	"bytes"
)

const (
	MaxPreAlloc = 1024 * 1024
)

var (
	ErrRangeIllegal = errors.New("range is illegal")
)

type Sds struct {
	// 空余的长度
	free int
	// 字符串长度
	len int
	// 实际字符串
	buf []byte
}

// 创建一个新的sds
func NewSds(str string) *Sds {
	sds := new(Sds)

	sds.len = utf8.RuneCountInString(str)
	sds.buf = []byte(str)
	sds.free = 0

	return sds
}

// 复制一个特定的sds
func DupSds(sds *Sds) *Sds {
	return NewSds(string(sds.buf))
}

// 比较两个sds
func CmpSds(s1 *Sds, s2 *Sds) int {
	return bytes.Compare(s1.buf, s2.buf)
}

// 返回sds长度
func (sds *Sds) SLen() int {
	return sds.len
}

// 返回可用空间
func (sds *Sds) Avail() int {
	return sds.free
}

// 扩展空间
func (sds *Sds) MakeRoomFor(addLen int) {
	if sds.Avail() > addLen {
		return
	}

	length := sds.SLen()
	newLen := length + addLen

	if newLen < MaxPreAlloc {
		newLen = newLen * 2
	} else {
		newLen = newLen + MaxPreAlloc
	}

	newBuf := make([]byte, len(sds.buf), newLen)
	copy(newBuf, sds.buf)

	sds.buf = newBuf
	sds.free = newLen - sds.SLen()
}

// 重置空字符串
func (sds *Sds) Clear() {
	sds.free = sds.SLen() + sds.free
	sds.len = 0
	sds.buf = []byte("")
}

// 拼接字符串到sds后面
func (sds *Sds) StrCat(str string) {
	addLen := utf8.RuneCountInString(str)

	sds.MakeRoomFor(addLen)

	sds.len += addLen
	sds.free = sds.free - addLen

	sds.buf = append(sds.buf, []byte(str)...)
}

// 把另一个sds字符串拼接到末尾
func (sds *Sds) StrCatSds(newSds Sds) {
	sds.StrCat(string(newSds.buf))
}

// 把一个字符串copy到sds
func (sds *Sds) Copy(str string) {
	addLen := utf8.RuneCountInString(str)

	total := sds.SLen() + sds.Avail()
	if total < addLen {
		sds.MakeRoomFor(addLen - total)
		total = sds.SLen() + sds.Avail()
	}

	sds.buf = []byte(str)
	sds.len = addLen
	sds.free = total - addLen
}

// 截取sds
func (sds *Sds) Range(start, end int) error {
	if start >= sds.SLen() || end > sds.SLen() {
		return ErrRangeIllegal
	}

	if end < 0 {
		end = sds.SLen() - int(math.Abs(float64(end))) + 1
	}

	if start < 0 {
		start = sds.SLen() - int(math.Abs(float64(start)))
	}

	var newBuf []byte
	if end > 0 {
		newBuf = sds.buf[start:end]
	} else {
		newBuf = sds.buf[start:]
	}

	newLen := len(newBuf)

	sds.free = sds.free + (sds.len - newLen)
	sds.buf = newBuf
	sds.len = len(sds.buf)

	return nil
}

func (sds *Sds) Trim(cutset string) {
	result := strings.Trim(string(sds.buf), cutset)
	cutLen := utf8.RuneCountInString(result)

	sds.buf = []byte(result)
	sds.free = sds.Avail() + (sds.SLen() - cutLen)
	sds.len = cutLen
}

func (sds *Sds) String() string {
	return string(sds.buf)
}

