package main

import (
	"fmt"

	"github.com/oandrew/ipod"
	"github.com/sirupsen/logrus"
)

func FrameLogEntry(e *logrus.Entry, frame []byte) *logrus.Entry {
	return e.WithFields(logrus.Fields{
		"len": len(frame),
	})
}

func CommandLogEntry(e *logrus.Entry, cmd *ipod.Command) *logrus.Entry {
	if cmd == nil {
		return e
	}
	return e.WithFields(logrus.Fields{
		"id":   cmd.ID,
		"trx":  cmd.Transaction,
		"type": fmt.Sprintf("%T", cmd.Payload),
	})
}
