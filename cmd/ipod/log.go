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

func PacketLogEntry(e *logrus.Entry, pkt *ipod.Packet) *logrus.Entry {
	if pkt == nil {
		return e
	}
	return e.WithFields(logrus.Fields{
		"id":   pkt.ID,
		"trx":  pkt.Transaction,
		"type": fmt.Sprintf("%T", pkt.Payload),
	})
}
