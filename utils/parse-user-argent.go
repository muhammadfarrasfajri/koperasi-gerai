package utils

import (
	"fmt"

	"github.com/mssola/user_agent"
)

func ParseDeviceInfo(uaRaw string) string {
	ua := user_agent.New(uaRaw)

	browser, version := ua.Browser()
	os := ua.OS()

	deviceType := "Desktop"
	if ua.Mobile() {
		deviceType = "Mobile"
	}

	return fmt.Sprintf(
		"%s %s / %s / %s",
		browser,
		version,
		os,
		deviceType,
	)
}
