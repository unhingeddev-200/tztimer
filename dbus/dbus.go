package dbus

import (
	"context"
	"fmt"

	bus "github.com/godbus/dbus/v5"
)

type Dbus struct {
	session         *bus.Conn
	notificationSvc bus.BusObject
}

func (d *Dbus) Notify(title string, body string, expireMili int) error {
	res := new(uint)
	err := d.notificationSvc.Call("org.freedesktop.Notifications.Notify", 0, "tztimer", uint(1), "", title, body, []string{}, map[string]any{}, expireMili).Store(res)
	if err != nil {
		return err
	}
	return nil
}

func NewDbus(ctx context.Context) (*Dbus, error) {
	b := new(Dbus)

	session, err := bus.ConnectSessionBus(bus.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}
	b.session = session
	b.notificationSvc = session.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	return b, nil
}
