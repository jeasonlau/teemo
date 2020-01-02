// +build linux freebsd netbsd openbsd

package beeep

// Alert displays a desktop notification and plays a beep.
func Alert(appID, title, message, appIcon string) error {
	if err := Notify(appID, title, message, appIcon); err != nil {
		return err
	}
	return Beep(DefaultFreq, DefaultDuration)
}
