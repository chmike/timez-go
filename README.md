# TimeZ 

Go package implementing Timez. A timez is a compact and convenient binary encoding
for UTC time in microseconds units and a local time offset.

## Why ? 

[ISO](https://en.wikipedia.org/wiki/International_Organization_for_Standardization)
issued the standard [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) to 
normalize text encoded time values. The [IETF](https://en.wikipedia.org/wiki/Internet_Engineering_Task_Force) issued an
equivalent standard ([RFC3339]()https://www.ietf.org/rfc/rfc3339.txt) for 
the Internet. 

These time encodings support combining a local time with a time offset to 
UTC. But they are text only. Timez is the binary equivalent of these time values.

A timez value encodes the UTC time with microsecond precision and the local time
offset in a 64 bit unsigned integer. Both values can be retrieved independently.
This makes timez a convenient stamp value in binary messages crossing local time
offset zones boundaries. 

The most interseting property of timez values is that comparing the integer 
values is the same as comparing by the timez values by UTC time and, when
equal, by time offset. Timez values are then convenient and efficient to use
as key in a sorted table, or as indexed value stored in a database.

## Timez properties

1. A timez can encode an absolute time or a time interval. 
2. Timez are encoded in a 64 bit unsigned integer to be simple and efficient to
mnipulate, stored in a database and be indexed ;
3. Timez provides time with micro second resolution which is acceptable for
Internet applications since a photon travels at most 300m in 1 micro second in
vacuum ;
4. The time offset is in minutes and covers the range -17:03 to 17:03 ;
5. The covered time range is from Jan, 1 1970 to approximatly xx, x 2255 ;
6. Comparing the integer encoding a timez yields the same result as comparing
the UTC time, and, when equal, comparing the time offset.

## Timez encoding

A Timez encodes the number of micro seconds elapsed since 
1970-01-01T00:00:00.000000Z in the 53 most significant bits of a 64 bit
unsigned integer.

The time offset is encoded in the 11 less significant bits as a number
of minutes relative to 1024. Thus the offset value 1024 is the time 
offset 00:00, the value 984 is the time offset -01:00, and the value 
1084 it the time offset +01:00. The time offset value 0 is invalid.

	64                                 11        0   bits
	|___________________  ______________|________|
	|__________________//_______________|________|
	|  number of micro seconds elapsed  |  time  |
    | since 1970-01-01T00:00:00:.000000 | offset |

The default initializer of timez values yields an invalid timez value. 

Note: the curent timez encoding of this package differ from the 
"github/chmike/timez" C library in that the epoch is different.

Feedback is welcome.

**Note:** This Work In Progress (WIP) and the encoding may change at any time.

