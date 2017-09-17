/*
Package timez provides compact binary encoded UTC time with time offset.

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
*/
package timez

import (
	"time"
)

// Time is a timez values.
type Time int64

// Invalid is an invalid time value.
const Invalid Time = 0

const offsetMask int64 = 0x7FF
const microsMask int64 = 0x3FF
const utcEpochSec int64 = int64(1)<<33 - int64(1)<<31
const minUtcSec int64 = -4503599627
const maxUtcSec int64 = 4503599627
const minOffset = -1023 * 60
const maxOffset = 1023 * 60

// Now return the current timez time.
func Now() Time {
	return From(time.Now())
}

// From converts a time.Time value to a timez.Time value.
// Return timez.Invalid if the UTC time or time offset are out of range.
// Valid UTC time are in the range year 1960 to year 2240. Valid time
// offset are in minute units and in the range -17:03 to 17:03.
func From(t time.Time) Time {
	var tz = t.Unix() - utcEpochSec
	if tz < minUtcSec || tz > maxUtcSec {
		return Invalid
	}
	_, offset := t.Zone()
	if offset < minOffset || offset > maxOffset || offset%60 != 0 {
		return Invalid
	}
	return Time(tz<<21 | int64(t.Nanosecond())/1000<<11 | int64(offset)/60 + 1024)
}

var fixedZones = map[uint16]*time.Location{}

func getLocation(o uint16) *time.Location {
	if l := fixedZones[o]; l != nil {
		return l
	}
	l := time.FixedZone("", (int(o)-1024)*60)
	fixedZones[o] = l
	return l
}

// Time converts a timez value to a time.Time value. The Location is set to a
// tine.FixedZone. The time offset will remain constant since the location is
// unknown and the daytime saving time change rule evolution is unknown.
// The returned time will be a zero time.Time value if the timez is invalid.
// This can be determined by calling the IsZero() method.
func (tz Time) Time() time.Time {
	var v = int64(tz)
	var offset = uint16(v & 0x7FF)
	if offset == 0 {
		return time.Time{}
	}
	var nsec = ((v >> 11) & 0x3FF) * 1000
	var sec = (v >> 21)
	return time.Unix(sec, nsec).In(getLocation(offset))
}

// String return the time in RFC3339 reprensentation with microseconds.
func (tz Time) String() string {
	return tz.Time().Format("2006-01-02T15:04:05.999999Z07:00")
}

// IsZero return true if the timez value is zero.
func (tz Time) IsZero() bool {
	return int64(tz)&offsetMask == 0
}
