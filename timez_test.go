package timez

import (
	"testing"
	"time"
)

func TestTimezTime(t *testing.T) {
	var tzNow = Now()
	var tzNowStr = tzNow.String()
	var tNow = tzNow.ToTime()
	var tNowStr = tNow.Format("2006-01-02T15:04:05.999999Z07:00")
	if tzNowStr != tNowStr {
		t.Errorf("got timez now %s, expected %s", tzNowStr, tNowStr)
	}

	// test bogus time offset values
	var tm = time.Unix(int64(1234567890), int64(0)).In(time.FixedZone("", 33))
	var tz = FromTime(tm)
	if tz != Invalid {
		t.Errorf("got timez %d, expected %d", uint64(tz), uint64(Invalid))
	}
	tm = time.Unix(int64(1234567890), int64(0)).In(time.FixedZone("", 17*3600+4*60))
	tz = FromTime(tm)
	if tz != Invalid {
		t.Errorf("got timez %d, expected %d", uint64(tz), uint64(Invalid))
	}
	tm = time.Unix(int64(1234567890), int64(0)).In(time.FixedZone("", -(17*3600 + 4*60)))
	tz = FromTime(tm)
	if tz != Invalid {
		t.Errorf("got timez %d, expected %d", uint64(tz), uint64(Invalid))
	}

	// test bogus sec values
	tm = time.Unix(-int64(1), int64(0)).In(time.FixedZone("", 0))
	tz = FromTime(tm)
	if tz != Invalid {
		t.Errorf("got timez %d, expected %d", uint64(tz), uint64(Invalid))
	}
	tm = time.Unix(int64(1<<53), int64(0)).In(time.FixedZone("", 0))
	tz = FromTime(tm)
	if tz != Invalid {
		t.Errorf("got timez %d, expected %d", uint64(tz), uint64(Invalid))
	}

	// test bogus timez value
	tz = FromUint64(123456789 << offsetBits)
	if tz.String() != "0001-01-01T00:00:00Z" {
		t.Errorf("got time %s, expected %s", tz.String(), "0001-01-01T00:00:00Z")
	}
	tz = Invalid
	if tz.String() != "0001-01-01T00:00:00Z" {
		t.Errorf("got time %s, expected %s", tz.String(), "0001-01-01T00:00:00Z")
	}
}

func TestTimezUint64(t *testing.T) {
	var val uint64 = 123456789<<offsetBits | 1084
	var tz = FromUint64(val)
	if tz.ToUint64() != val {
		t.Errorf("got %v, expected %v", tz.ToUint64(), val)
	}
	tz = FromUint64(123456789 << offsetBits)
	if tz != Invalid {
		t.Errorf("got timez %v, expected %v", uint64(tz), uint64(Invalid))
	}

}

func TestTimezOffset(t *testing.T) {
	var t0 = FromUint64(123456789<<offsetBits | 1084)
	if t0.Offset() != 3600 {
		t.Errorf("got offset %d, expected %d", t0.Offset(), 3600)
	}

	if err := t0.SetOffset(-3600); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if t0.Offset() != -3600 {
		t.Errorf("got offset %d, expected %d", t0.Offset(), -3600)
	}

	if err := t0.SetOffset(33); err == nil {
		t.Errorf("unexpected nil error")
	}
	if err := t0.SetOffset(-(17*3600 + 4*60)); err == nil {
		t.Errorf("unexpected nil error")
	}
	if err := t0.SetOffset(17*3600 + 4*60); err == nil {
		t.Errorf("unexpected nil error")
	}
}
