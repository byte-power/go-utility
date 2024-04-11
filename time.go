package utility

import (
	"strings"
	"time"
)

const timeLocationName = "Asia/Shanghai"
const TimeLayout = "2006-01-02 15:04:05"

var (
	CNTimeLocation *time.Location
)

const (
	CNGMT = "+0800"
	// Time duration for one day.
	TimeDay = 24 * time.Hour
)

func GetCNCurrentTime() time.Time {
	return time.Now().In(CNTimeLocation)
}

func TimeToCNFormat(t time.Time) string {
	return t.In(CNTimeLocation).Format(TimeLayout)
}

func CNFormatToTime(timeString string) time.Time {
	t, err := time.ParseInLocation(TimeLayout, timeString, CNTimeLocation)
	if err != nil {
		return time.Time{}
	}
	return t
}

func GetCNFormatCurrentTime() string {
	return time.Now().In(CNTimeLocation).Format(TimeLayout)
}

func DayDiff(date1, date2 time.Time, loc *time.Location) int {
	y1, m1, d1 := date1.In(loc).Date()
	y2, m2, d2 := date2.In(loc).Date()
	date1 = time.Date(y1, m1, d1, 0, 0, 0, 0, loc)
	date2 = time.Date(y2, m2, d2, 0, 0, 0, 0, loc)
	return int(date1.Sub(date2) / TimeDay)
}

func IsDateEqual(date1, date2 time.Time, loc *time.Location) bool {
	y1, m1, d1 := date1.In(loc).Date()
	y2, m2, d2 := date2.In(loc).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func MillisecondOfTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

type TimeInDay struct {
	Hour   int
	Minute int
	Second int
}

func MakeTimeInDay(hour, minute, second int) TimeInDay {
	return TimeInDay{Hour: hour, Minute: minute, Second: second}
}

func (t TimeInDay) IsValid() bool {
	return t.Hour >= 0 && t.Hour <= 24 &&
		t.Minute >= 0 && t.Minute <= 60 &&
		t.Second >= 0 && t.Second <= 60
}

func (t TimeInDay) TimeOnDay(day time.Time) time.Time {
	y, m, d := day.Date()
	return time.Date(y, m, d, t.Hour, t.Minute, t.Second, 0, day.Location())
}

// Compare with the other time in the day of the other time.
//   - Return: receiver -1 early, 0 same, 1 later than the other time
func (t TimeInDay) Compare(other time.Time) int8 {
	_t := t.TimeOnDay(other)
	_ts := _t.Unix()
	ots := other.Unix()
	if _ts == ots {
		return 0
	} else if _ts > ots {
		return 1
	}
	return -1
}

// ParseTimeInDay would parse TimeInDay from time string, e.g. 09 or 8:25 or 22:03:2.
func ParseTimeInDay(t string) (inst TimeInDay) {
	if t == "" {
		inst.Hour = -1
		return
	}
	comps := strings.Split(t, ":")
	l := len(comps)
	if l == 0 {
		inst.Hour = -1
		return
	}
	inst.Hour = int(AnyToInt64(comps[0]))
	if l > 1 {
		inst.Minute = int(AnyToInt64(comps[1]))
		if l > 2 {
			inst.Second = int(AnyToInt64(comps[2]))
		}
	}
	return
}
