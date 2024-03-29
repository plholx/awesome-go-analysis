package main

import (
	"fmt"
	"time"
)

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

// Delay 延迟阻塞至某一时间
// isTomorrow 参数表示指定时间是否是明天的时刻
//
// 如果指定的时间为今天，且时间小于当前时间则不会延迟
func Delay(hour, min, sec, nsec int, isTomorrow bool) {
	targetDate := time.Now()
	if isTomorrow {
		targetDate = targetDate.Add(time.Hour * 24)
	}
	targetTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), hour, min, sec, nsec, targetDate.Location())
	duration := targetTime.Sub(time.Now())
	if duration > 0 {
		<-time.NewTimer(duration).C
	}
}

// TimeSince 转换指定时间为便于识别的形式
// 代码复制自 https://github.com/gogs/gogs/blob/436dd6c0a4549e09069193879046a2a202d6291e/pkg/tool/tool.go
func TimeSince(then *time.Time) string {
	now := time.Now()

	lbl := "ago"
	diff := now.Unix() - then.Unix()
	if then.After(now) {
		lbl = "from now"
		diff = then.Unix() - now.Unix()
	}

	switch {
	case diff <= 0:
		return "now"
	case diff <= 2:
		return fmt.Sprintf("1 second %s", lbl)
	case diff < 1*Minute:
		return fmt.Sprintf("%d seconds %s", diff, lbl)
	case diff < 2*Minute:
		return fmt.Sprintf("1 minute %s", lbl)
	case diff < 1*Hour:
		return fmt.Sprintf("%d minutes %s", diff/Minute, lbl)
	case diff < 2*Hour:
		return fmt.Sprintf("1 hour %s", lbl)
	case diff < 1*Day:
		return fmt.Sprintf("%d hours %s", diff/Hour, lbl)
	case diff < 2*Day:
		return fmt.Sprintf("1 day %s", lbl)
	case diff < 1*Week:
		return fmt.Sprintf("%d days %s", diff/Day, lbl)
	case diff < 2*Week:
		return fmt.Sprintf("1 week %s", lbl)
	case diff < 1*Month:
		return fmt.Sprintf("%d weeks %s", diff/Week, lbl)
	case diff < 2*Month:
		return fmt.Sprintf("1 month %s", lbl)
	case diff < 1*Year:
		return fmt.Sprintf("%d months %s", diff/Month, lbl)
	case diff < 2*Year:
		return fmt.Sprintf("1 year %s", lbl)
	default:
		return fmt.Sprintf("%d years %s", diff/Year, lbl)
	}
}