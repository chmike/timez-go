# TimeZ 

Package timez provides compact binary encoded UTC time with time offset.

This a Work In Progress (WIP).

By giving access to the UTC time and time offest, it is possible to
convert the timez value to different earth time referentials.

Typical use is for time stamp stored in a database containing data generated
in earth locations with different time offset, or in messages transmitted
between locations with different time offset.

A timez is a 64bit signed integer which is compact and can easily be stored in
existing databases. Sorting these integers will sort them by UTC time, and
time offsets when the UTC time are equal.

A timez split the 64 bits of an int64 into the parts. The UTC time is stored
as the number of microsecond offsets since an epoch in the 53 most significant
bits as a two's complement integer value. The time offset is stored as the
minutes elapsed since 00:00 plus 1024 minutes. For instance the value 1024
is the offset 00:00, 1084 the offset +01:00, and 604 the offset -07:00.

	64                31             11        0   bits
	|_________  _______|______________|________|
	|________//________|______________|________|
	| seconds relative | microseconds |  time  |
    | to an UTC epoch  |              | offset |

The epoch is picked so that the unix time period is covered and beyond. The
smallest representable time value is Jan 1 1970 00:00:00 minus 2^32-1 seconds
as for the unix time. Since we have 33 bits to encode the number of seconds
since that time, the timez epoch is 2^33 - 2^31. The biggest time value that 
can be represented with a timez value is in the year 2174.

Note: the curent timez encoding in this package is different from the 
"github/chmike/timez" C library. The epoch is different and this package uses
a bit field to store microseconds. It is expected to more efficient to convert 
with the Go time representation, at the detriment of the time range it can 
cover. An evaluation of the trade off is still required. 

Feedback is welcome.