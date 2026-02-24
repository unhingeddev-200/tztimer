package timer

import (
	"fmt"
	"time"

	"tztimer/dbus"

	"github.com/sirupsen/logrus"
)

type Timer struct {
	tz       *time.Location
	notifyAt time.Time
	bus      *dbus.Dbus
}

func New(bus *dbus.Dbus, tz *time.Location, notifyAt string) (*Timer, error) {
	timeOnly, err := time.ParseInLocation(time.TimeOnly, notifyAt, tz)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}
	now := time.Now().In(tz)
	target := time.Date(now.Year(), now.Month(), now.Day(), timeOnly.Hour(), timeOnly.Minute(), timeOnly.Second(), timeOnly.Nanosecond(), tz)

	return &Timer{
		tz:       tz,
		notifyAt: target,
		bus:      bus,
	}, nil
}

func (t *Timer) Start() <-chan time.Time {
	var d time.Duration
	logrus.Debug(time.Now().In(t.notifyAt.Location()))
	logrus.Debug(t.notifyAt)
	if time.Now().In(t.notifyAt.Location()).After(t.notifyAt) {
		d = 0
	} else {
		d = t.notifyAt.Sub(time.Now().In(t.notifyAt.Location()))
		t.bus.Notify(fmt.Sprintf("Alarm set for: %s", d), "", 5000)
	}
	logrus.WithField("TimeDiff", d).Debug("Time - Now")
	return time.NewTimer(d).C
}
