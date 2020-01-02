// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/tadvi/systray"
	"golang.org/x/sys/windows/registry"
	toast "gopkg.in/toast.v1"
)

var isWindows10 bool
var applicationID string

func init() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer k.Close()

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return
	}

	isWindows10 = maj == 10

	if isWindows10 {
		applicationID = appID()
	}
}

// Notify sends desktop notification.
func Notify(appID, title, message, appIcon string) error {
	if isWindows10 {
		return toastNotify(appID, title, message, appIcon)
	}

	err := baloonNotify(title, message, appIcon, false)
	if err != nil {
		e := msgNotify(title, message)
		if e != nil {
			return errors.New("beeep: " + err.Error() + "; " + e.Error())
		}
	}

	return nil

}

func msgNotify(title, message string) error {
	msg, err := exec.LookPath("msg")
	if err != nil {
		return err
	}
	cmd := exec.Command(msg, "*", "/TIME:3", title+"\n\n"+message)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func baloonNotify(title, message, appIcon string, bigIcon bool) error {
	tray, err := systray.New()
	if err != nil {
		return err
	}

	err = tray.ShowCustom(pathAbs(appIcon), title)
	if err != nil {
		return err
	}

	go func() {
		tray.Run()
		time.Sleep(3 * time.Second)
		tray.Stop()
	}()

	return tray.ShowMessage(title, message, bigIcon)
}

func toastNotify(appID, title, message, appIcon string) error {
	notification := toastNotification(appID, title, message, pathAbs(appIcon))
	return notification.Push()
}

func toastNotification(appID, title, message, appIcon string) toast.Notification {
	if len(appIcon) > 0 {
		return toast.Notification{
			AppID:   appID,
			Title:   title,
			Message: message,
			Icon:    appIcon,
		}
	}
	return toast.Notification{
		AppID:   appID,
		Title:   title,
		Message: message,
	}
}

func appID() string {
	defID := "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd := exec.Command("powershell", "Get-StartApps")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return defID
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "powershell.exe") {
			sp := strings.Split(line, " ")
			if len(sp) > 0 {
				return sp[len(sp)-1]
			}
		}
	}

	return defID
}
