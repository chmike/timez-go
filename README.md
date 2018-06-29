[![GoDoc](https://godoc.org/github.com/chmike/timez-go?status.svg)](https://godoc.org/github.com/chmike/timez-go)
[![Build](https://travis-ci.org/chmike/timez-go.svg?branch=master)](https://travis-ci.org/chmike/timez-go?branch=master)
[![Coverage](https://coveralls.io/repos/github/chmike/timez-go/badge.svg?branch=master)](https://coveralls.io/github/chmike/timez-go?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/chmike/timez-go)](https://goreportcard.com/report/github.com/chmike/timez-go)
![Status](https://img.shields.io/badge/status-beta-orange.svg)
![release](https://img.shields.io/github/release/chmike/timez-go.svg)

# TimeZ 

Go package implementing Timez. A timez is 64 bit unsigned integer combining an 
UTC time value in microsecond, and a time zone time offset in minutes. Comparing 
two such unisgned integer yields the same as if the UTC time was first compared, 
and if equal, the time zone offset were compared.

## Why ? 

[ISO](https://en.wikipedia.org/wiki/International_Organization_for_Standardization)
issued the standard [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) to 
normalize text encoded time values. The [IETF](https://en.wikipedia.org/wiki/Internet_Engineering_Task_Force) issued an
equivalent standard ([RFC3339]()https://www.ietf.org/rfc/rfc3339.txt) for 
the Internet. These time encodings combine a local time with a time offset relative to 
UTC. But they are text only. Timez is the binary equivalent of these time values.

## What is a timez ?

A timez value encodes the UTC time with microsecond precision and the local time
offset in a 64 bit unsigned integer. Both values can be retrieved independently.
This makes timez a convenient stamp value in binary messages crossing time zone
boundaries. 

The most interseting property of timez values is that comparing the integer 
values is equivalent to compare UTC time, and when equal, to compare time zone
time offsets. Timez values are then convenient and efficient to use
as key in a sorted table, or as ordered index key in a database.

## The microsecond resolution

The microsecond UTC time resolution is a compromise. A nanosecond resolution
would have been preferable, but it wouldn't fit in a 64bit integer with the
time offset.
[NTP](https://en.wikipedia.org/wiki/Network_Time_Protocol) can at very best
synchronize around a few tens of microseconds. With GPS, the  best time 
synchronization we could get is around a few tens of nanoseconds. Since
a photon can travel at most 300m in a microsecond in vacuum, for message
stamps with Internet application, a microsecond precision is an acceptable
compromise.

## Timez encoding

A Timez encodes the number of microseconds elapsed since 
1970-01-01T00:00:00.000000Z in the 53 most significant bits of a 64 bit
unsigned integer. The 53 bit microsecond counter will wrap 285 years after
1970, thus in the year 2255.

The time offset is encoded in the 11 less significant bits as a number
of minutes relative to 1024. Thus the offset value 1024 is the time 
offset 00:00, the value 984 is the time offset -01:00, and the value 
1084 it the time offset +01:00. The time offset value 0 is invalid.

    64                                11        0   bits
    |__________________  ______________|________|
    |_________________//_______________|________|
    |  number of microseconds elapsed  |  time  |
    | since 1970-01-01T00:00:00.000000 | offset |

The default initializer of Go timez values yields an invalid timez value. 

Note: the curent timez encoding of this package differ from the 
"github/chmike/timez" C library in that we now use 64 bit unsigned integer.

Feedback is welcome.

**Note:** This Work In Progress (WIP) and the encoding may change at any time.

