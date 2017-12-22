package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

type TextFormatter struct {
	DisableColors bool
	isTerminal    bool
	sync.Once
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	f.print(b, entry, keys)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) colored() bool {
	return f.isTerminal && !f.DisableColors
}

func (f *TextFormatter) print(b *bytes.Buffer, entry *logrus.Entry, keys []string) {

	levelText := strings.ToUpper(entry.Level.String())[0:4]
	ts := entry.Time.Sub(baseTimestamp)
	tsText := fmt.Sprintf("%04d.%06d", ts/time.Second, (ts%time.Second)/time.Microsecond)
	if f.colored() {
		var levelColor int
		switch entry.Level {
		case logrus.DebugLevel:
			levelColor = gray
		case logrus.WarnLevel:
			levelColor = yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = red
		default:
			levelColor = blue
		}
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, tsText, entry.Message)
		for _, k := range keys {
			v := entry.Data[k]
			fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%v", levelColor, k, v)
		}
	} else {
		fmt.Fprintf(b, "%s[%s] %-44s ", levelText, tsText, entry.Message)
		for _, k := range keys {
			v := entry.Data[k]
			fmt.Fprintf(b, " %s=%v", k, v)
		}
	}

}
