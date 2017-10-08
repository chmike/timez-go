/*
Package timez provides compact binary encoding of UTC time and time offset.

A timez value encodes the UTC time with microsecond precision and the local time
offset in a 64 bit unsigned integer. Both values can be retrieved independently.
This makes timez a convenient stamp value in binary messages crossing local time
offset zones boundaries.

The most interseting property of timez values is that comparing the integer
values is the same as comparing by the timez values by UTC time and, when
equal, by time offset. Timez values are then convenient and efficient to use
as key in a sorted table, or as indexed value stored in a database.

Timez encoding

A Timez encodes the number of micro seconds elapsed since
1970-01-01T00:00:00.000000Z in the 53 most significant bits of a 64 bit
unsigned integer.

The time offset is encoded in the 11 less significant bits as a number
of minutes relative to 1024. Thus the offset value 1024 is the time
offset 00:00, the value 984 is the time offset -01:00, and the value
1084 it the time offset +01:00.

	64                                11        0   bits
	|__________________  ______________|________|
	|_________________//_______________|________|
	|  number of microseconds elapsed  |  time  |
    | since 1970-01-01T00:00:00.000000 | offset |


The epoch is picked so that the unix time period is covered and beyond. The
smallest representable time value is Jan 1 1970 00:00:00 minus 2^32-1 seconds
as for the unix time. Since we have 33 bits to encode the number of seconds
since that time, the timez epoch is 2^33 - 2^31. The biggest time value that
can be represented with a timez value is in the year 2174.
*/
package timez

import (
	"errors"
	"sync"
	"time"
)

// Time is a timez values.
type Time uint64

// Invalid is an invalid time value.
const Invalid Time = 0
const offsetBits byte = 11
const offsetMask uint64 = 0x7FF
const utcEpochSec int64 = 0
const minUtcSec int64 = 0
const maxUtcSec int64 = ((1 << 53) / 1000000) - 1
const minOffsetSec = -1023 * 60
const maxOffsetSec = 1023 * 60
const utcOffset uint16 = 1024

// Now return the current timez time.
func Now() Time {
	return FromTime(time.Now())
}

// FromTime converts a time.Time value to a timez.Time value.
// Return timez.Invalid if the time.Time value can't be converted to a
// timez.Time value.
func FromTime(t time.Time) Time {
	var tsec = t.Unix() - utcEpochSec
	if tsec < minUtcSec || tsec > maxUtcSec {
		return Invalid
	}
	_, offset := t.Zone()
	if offset < minOffsetSec || offset > maxOffsetSec || offset%60 != 0 {
		return Invalid
	}
	var microsec = (uint64(tsec)*1000000 + uint64(t.Nanosecond())/1000)
	return Time(microsec<<offsetBits | uint64(offset)/60 + 1024)
}

var locationMap sync.Map // map[uint16]*time.Location{}

func getLocation(offset uint16) *time.Location {
	if l, ok := locationMap.Load(offset); ok {
		return l.(*time.Location)
	}

	l := time.FixedZone("", (int(offset)-1024)*60)
	locationMap.Store(offset, l)
	return l
}

// ToTime converts a timez value to a time.Time value. The Location is set to a
// time.FixedZone. The time offset will remain constant since the location is
// unknown and the daytime saving time change rule evolution is unknown.
// The returned time will be a zero time.Time value if the timez is invalid.
// This can be determined by calling the IsZero() method.
func (tz Time) ToTime() time.Time {
	var t = uint64(tz)
	var offset = uint16(t & offsetMask)
	if offset == 0 {
		return time.Time{}
	}
	t >>= offsetBits
	var nsec = int64(1000 * (t % 1000000))
	var sec = int64(t / 1000000)
	return time.Unix(sec, nsec).In(getLocation(offset))
}

// String return the time in RFC3339 reprensentation with microseconds.
func (tz Time) String() string {
	return tz.ToTime().Format("2006-01-02T15:04:05.999999Z07:00")
}

// ToUint64 return the timez value as a uint64 value.
func (tz Time) ToUint64() uint64 {
	return uint64(tz)
}

// FromUint64 return the uint64 value as a timez value. Return Invalid if
// the uint64 value is not a valid timez value.
func FromUint64(tz uint64) Time {
	if (tz & offsetMask) == 0 {
		return Invalid
	}
	return Time(tz)
}

// Offset return the local time offset in seconds.
func (tz Time) Offset() int {
	return (int(uint64(tz)&offsetMask) - 1024) * 60
}

// SetOffset sets the timez time offset to offset.
func (tz *Time) SetOffset(offset int) error {
	if offset < minOffsetSec || offset > maxOffsetSec || offset%60 != 0 {
		return errors.New("invalid timez time offset")
	}
	var t = uint64(*tz) & ^offsetMask
	*tz = Time(t | uint64((offset/60)+1024))
	return nil
}
