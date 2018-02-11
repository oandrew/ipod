/*
ipod talks to an ipod accessory using iap protocol.

# with debug logging
./ipod -d serve /dev/iap0

# save a trace file
./ipod -d serve -w ipod.trace /dev/iap0

# simulate incoming requests from a trace file
./ipod -d replay ./ipod.trace

# view a trace file
./ipod -d view ./ipod.trace


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


*/
package main
