package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/oandrew/ipod"
	"github.com/oandrew/ipod/hid"
	"github.com/oandrew/ipod/lingo-audio"
	"github.com/oandrew/ipod/lingo-dispremote"
	"github.com/oandrew/ipod/lingo-extremote"
	"github.com/oandrew/ipod/lingo-general"
	"github.com/oandrew/ipod/lingo-simpleremote"
	"github.com/oandrew/ipod/trace"
)

var log = logrus.StandardLogger()

func openDevice(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	stat, _ := f.Stat()
	if stat.Mode()&os.ModeCharDevice != os.ModeCharDevice {
		return nil, fmt.Errorf("not a char device")
	}
	return f, nil
}

func openTraceFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func newTraceFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
}

type UsageError struct {
	error
}

var logBufW *bufio.Writer

func main() {
	logOut := os.Stderr
	log.Formatter = &TextFormatter{
		Colored: checkIfTerminal(logOut),
	}

	logBufW = bufio.NewWriter(logOut)
	defer logBufW.Flush()

	log.Out = logBufW

	app := cli.NewApp()
	app.Name = "ipod"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Andrew Onyshchuk",
		},
	}
	app.Usage = "ipod accessory protocol client"
	app.HideVersion = true

	app.ErrWriter = os.Stderr
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "verbose logging",
		},
	}

	exitErrorHandler := func(c *cli.Context, err error) {
		if err != nil {
			if _, ok := err.(UsageError); ok {
				fmt.Fprintf(c.App.ErrWriter, "usage error: %v\n\n", err)
				cli.ShowCommandHelp(c, c.Command.Name)
			} else {
				fmt.Fprintf(c.App.ErrWriter, "error: %v\n\n", err)
			}
			os.Exit(1)
		}
	}
	app.ExitErrHandler = exitErrorHandler

	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			log.SetLevel(logrus.DebugLevel)
		}

		log.Debugf("Registered lingos:\n%s", ipod.DumpLingos())

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "serve",
			Aliases:   []string{"s"},
			ArgsUsage: "<dev>",
			Usage:     "respond to requests from a char device i.e. /dev/iap0",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "write-trace, w",
					Usage: "Write trace to a `file`",
				},
			},
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				if path == "" {
					return UsageError{fmt.Errorf("device path is missing")}
				}
				f, err := openDevice(path)
				le := log.WithField("path", path)
				if err != nil {
					le.WithError(err).Errorf("could not open the device")
					return err
				}
				le.Info("device opened")

				var rw io.ReadWriter = f
				if tracePath := c.String("write-trace"); tracePath != "" {
					traceFile, err := newTraceFile(tracePath)
					le := log.WithField("path", tracePath)
					if err != nil {
						le.WithError(err).Errorf("could not create a trace file")
						return err
					}
					le.Warningf("writing trace")
					rw = trace.NewTracer(traceFile, f)
				}

				reportR, reportW := hid.NewReportReader(rw), hid.NewReportWriter(rw)
				frameTransport := hid.NewTransport(reportR, reportW, hid.DefaultReportDefs)
				processFrames(frameTransport)
				return nil
			},
		},
		{
			Name:    "replay",
			Aliases: []string{"r"},
			Usage:   "respond to requests from a trace file",
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				if path == "" {
					return UsageError{cli.NewExitError("trace file path is missing", 1)}
				}

				f, err := openTraceFile(path)
				le := log.WithField("path", path)
				if err != nil {
					le.WithError(err).Errorf("could not open the trace file")
					return err
				}
				le.Warningf("trace file opened")

				tr := trace.NewReader(f)
				tdr := trace.NewTraceDirReader(tr, trace.DirIn)
				reportR, reportW := hid.NewReportReader(tdr), hid.NewReportWriter(ioutil.Discard)
				frameTransport := hid.NewTransport(reportR, reportW, hid.DefaultReportDefs)
				processFrames(frameTransport)
				return nil
			},
		},
		{
			Name:    "view",
			Aliases: []string{"v"},
			Usage:   "view a trace file",
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				if path == "" {
					return UsageError{cli.NewExitError("trace file path is missing", 1)}
				}

				f, err := openTraceFile(path)
				le := log.WithField("path", path)
				if err != nil {
					le.WithError(err).Errorf("could not open the trace file")
					return err
				}
				le.Warningf("trace file opened")
				tr := trace.NewReader(f)
				dumpTrace(tr)
				return nil
			},
		},
	}

	app.Setup()
	app.Run(os.Args)

}

func logFrame(frame []byte, err error, msg string) {
	le := FrameLogEntry(logrus.NewEntry(log), frame)
	if err != nil {
		le.WithError(err).Errorf(msg)
		return
	}
	le.Infof(msg)
	if log.Level == logrus.DebugLevel {
		spew.Fdump(logBufW, frame)
	}

}

func logPacket(pkt []byte, err error, msg string) {
	//le := PacketLogEntry(logrus.NewEntry(log), frame)
	le := log.WithField("len", len(pkt))
	if err != nil {
		le.WithError(err).Errorf(msg)
		return
	}
	le.Infof(msg)
	if log.Level == logrus.DebugLevel {
		spew.Fdump(logBufW, pkt)
	}
}

func logCmd(cmd *ipod.Command, err error, msg string) {
	le := CommandLogEntry(logrus.NewEntry(log), cmd)
	if err != nil {
		le.WithError(err).Errorf(msg)
		return
	}
	le.Infof(msg)
	if log.Level == logrus.DebugLevel {
		spew.Fdump(logBufW, cmd)
	}

}

func processFrames(frameTransport ipod.FrameReadWriter) {
	outFrameBuf := bytes.Buffer{}
	defer logBufW.Flush()
	for {
		logBufW.Flush()
		inFrame, err := frameTransport.ReadFrame()
		if err == io.EOF {
			break
		}
		logFrame(inFrame, err, "<< FRAME")
		if err != nil {
			continue
		}

		packetReader := ipod.NewPacketReader(bytes.NewReader(inFrame))
		inCmdBuf := ipod.CmdBuffer{}
		for {
			inPacket, err := packetReader.ReadPacket()
			if err == io.EOF {
				break
			}
			logPacket(inPacket, err, "<< PACKET")
			if err != nil {
				continue
			}

			var inCmd ipod.Command
			inCmdErr := inCmd.UnmarshalBinary(inPacket)
			logCmd(&inCmd, inCmdErr, "<< CMD")
			inCmdBuf.WriteCommand(&inCmd)
		}

		outCmdBuf := ipod.CmdBuffer{}
		for i := range inCmdBuf.Commands {
			//todo: check return error
			handlePacket(&outCmdBuf, inCmdBuf.Commands[i])
		}

		outFrameBuf.Reset()
		packetWriter := ipod.NewPacketWriter(&outFrameBuf)
		for i := range outCmdBuf.Commands {
			outCmd := outCmdBuf.Commands[i]
			logCmd(outCmd, nil, ">> CMD")

			outPacket, err := outCmd.MarshalBinary()
			logPacket(outPacket, err, ">> PACKET")

			packetWriter.WritePacket(outPacket)
		}

		if outFrameBuf.Len() == 0 {
			continue
		}
		outFrame := outFrameBuf.Bytes()
		outFrameErr := frameTransport.WriteFrame(outFrame)
		logFrame(outFrame, outFrameErr, ">> FRAME")

	}
	log.Warnf("EOF")
}

var devGeneral = &DevGeneral{}

func handlePacket(cmdWriter ipod.CommandWriter, cmd *ipod.Command) {
	switch cmd.ID.LingoID() {
	case general.LingoGeneralID:
		general.HandleGeneral(cmd, cmdWriter, devGeneral)
		if _, ok := cmd.Payload.(*general.RetDevAuthenticationSignature); ok {
			audio.Start(cmdWriter)
		}
	case simpleremote.LingoSimpleRemotelID:
		//todo
		log.Warn("Lingo SimpleRemote is not supported yet")
	case dispremote.LingoDisplayRemoteID:
		dispremote.HandleDispRemote(cmd, cmdWriter, nil)
	case extremote.LingoExtRemotelID:
		extremote.HandleExtRemote(cmd, cmdWriter, nil)
	case audio.LingoAudioID:
		audio.HandleAudio(cmd, cmdWriter, nil)
	}
}
func dirPrefix(dir trace.Dir, text string) string {
	switch dir {
	case trace.DirIn:
		return "<< " + text
	case trace.DirOut:
		return ">> " + text
	default:
		return "?? " + text
	}
}
func dumpTrace(tr *trace.Reader) {
	q := trace.Queue{}
	for {
		var msg trace.Msg
		err := tr.ReadMsg(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		q.Enqueue(&msg)
	}

	for {
		head := q.Head()
		if head == nil {
			break
		}
		dir := head.Dir
		tdr := trace.NewQueueDirReader(&q, dir)
		d := hid.NewDecoderDefault(hid.NewReportReader(tdr))

		frame, err := d.ReadFrame()
		if err == io.EOF {
			break
		}
		logFrame(frame, err, dirPrefix(dir, "FRAME"))
		if err != nil {
			continue
		}

		packetReader := ipod.NewPacketReader(bytes.NewReader(frame))
		for {
			packet, err := packetReader.ReadPacket()
			if err == io.EOF {
				break
			}
			logPacket(packet, err, dirPrefix(dir, "PACKET"))
			if err != nil {
				continue
			}

			var cmd ipod.Command
			cmdErr := cmd.UnmarshalBinary(packet)
			logCmd(&cmd, cmdErr, dirPrefix(dir, "CMD"))
		}
	}
	log.Warnf("EOF")
}
