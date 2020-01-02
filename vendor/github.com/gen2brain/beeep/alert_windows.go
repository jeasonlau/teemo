// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	toast "gopkg.in/toast.v1"
)

// Alert displays a desktop notification and plays a default system sound.
func Alert(appID, title, message, appIcon string) error {
	if isWindows10 {
		note := toastNotification(appID, title, message, pathAbs(appIcon))
		note.Audio = toast.Default
		return note.Push()
	}

	if err := Notify(appID, title, message, appIcon); err != nil {
		return err
	}
	return Beep(DefaultFreq, DefaultDuration)
}
