package godog

import (
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog/colors"
	"github.com/cucumber/messages-go/v10"
)

var (
	red    = colors.Red
	redb   = colors.Bold(colors.Red)
	green  = colors.Green
	blackb = colors.Bold(colors.Black)
	yellow = colors.Yellow
	cyan   = colors.Cyan
	cyanb  = colors.Bold(colors.Cyan)
	whiteb = colors.Bold(colors.White)
)

// repeats a space n times
func s(n int) string {
	if n < 0 {
		n = 1
	}
	return strings.Repeat(" ", n)
}

var timeNowFunc = func() time.Time {
	return time.Now()
}

func trimAllLines(s string) string {
	var lines []string
	for _, ln := range strings.Split(strings.TrimSpace(s), "\n") {
		lines = append(lines, strings.TrimSpace(ln))
	}
	return strings.Join(lines, "\n")
}

type sortPicklesByID []*messages.Pickle

func (s sortPicklesByID) Len() int { return len(s) }
func (s sortPicklesByID) Less(i, j int) bool {
	iID := mustConvertStringToInt(s[i].Id)
	jID := mustConvertStringToInt(s[j].Id)
	return iID < jID
}
func (s sortPicklesByID) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func mustConvertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}
