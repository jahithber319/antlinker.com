// timeFormat
package tools

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var _ = fmt.Println

var weekdays = [...]string{
	"星期日",
	"星期一",
	"星期二",
	"星期三",
	"星期四",
	"星期五",
	"星期六",
}

// 获取当前时间戳
func GetTimeStamp() int64 {
	return time.Now().Unix()
}

// 获取当前时间戳（带纳秒）
func GetTimeStampNamo() int64 {
	return time.Now().UnixNano()
}

// 按照格式获取当前日期，例如：yyyy-MM-dd HH:mm:ss
func GetDateWithFormat(format string) (string, error) {
	f := ""
	if strings.Contains(format, "y") {
		if strings.Count(format, "yyyy") == 1 && strings.Count(format, "y") == 4 {
			f = strings.Replace(format, "yyyy", "2006", -1)
		} else {
			return "", errors.New("日期格式错误")
		}
	}

	if strings.Contains(format, "M") {
		if strings.Count(format, "MM") == 1 && strings.Count(format, "M") == 2 {
			f = strings.Replace(f, "MM", "01", -1)
		} else {
			return "", errors.New("日期格式错误")
		}
	}

	if strings.Contains(format, "d") {
		if strings.Count(format, "dd") == 1 && strings.Count(format, "d") == 2 {
			f = strings.Replace(f, "dd", "02", -1)
		} else {
			return "", errors.New("日期格式错误")
		}
	}

	if strings.Contains(format, "H") {
		if strings.Count(format, "HH") == 1 && strings.Count(format, "H") == 2 {
			f = strings.Replace(f, "HH", "15", -1)
		} else {
			return "", errors.New("日期格式错误")
		}
	}

	if strings.Contains(format, "m") {
		if strings.Count(format, "mm") == 1 && strings.Count(format, "m") == 2 {
			f = strings.Replace(f, "mm", "04", -1)
		} else {
			return "", errors.New("日期格式错误")
		}
	}

	if strings.Contains(format, "s") {
		if strings.Count(format, "ss") == 1 && strings.Count(format, "s") == 2 {
			f = strings.Replace(f, "ss", "05", -1)
		} else {
			return "", errors.New("日期格式错误")
		}
	}

	return time.Now().Format(f), nil
}

// 返回给出日期是周几，英文
func GetWeekdayEn(t time.Time) string {
	return t.Weekday().String()
}

// 返回给出日期是周几，中文
func GetWeekdayZh(t time.Time) string {
	return weekdays[t.Weekday()]
}
