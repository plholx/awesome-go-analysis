package main_test

import (
	"fmt"
	analysis "github.com/plholx/awesome-go-analysis"
	"testing"
	"time"
)

func TestTimeSince(t *testing.T) {
	now := time.Now()
	timeStr := analysis.TimeSince(&now)
	fmt.Println(timeStr)
}