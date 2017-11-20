/*
ipod talks to an ipod accessory using iap protocol.

It has two modes:
Using -d flag, it will read and write iap packets
from/to specified char device i.e. the one provided by ipod-gadget.
Using -r flag, it will read incoming iap packets
from  a trace file and discard responses. Used for testing

With the -w flag, it will also save the trace
to a specified file that can be later used with -r

Each line of a trace file starts with a
 '< ' for incoming requests
 '> ' for outgoing responses
followed by the hex dump of the data
i.e. a trace file

 < 00 01 02
 > 02 01 00

represents an incoming request byte sequence from the accessory
 0x00,0x01,0x02

and an outgoing response byte sequence from the ipod
 0x02,0x01,0x00


The -v flag enables verbose logging including detailed structure
on each request/response

*/
package main
