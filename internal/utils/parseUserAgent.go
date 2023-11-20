package utils

import "github.com/mileusna/useragent"

type UserAgent struct {
	Device  string
	OS      string
	Browser string
}

func ParseUserAgent(userAgent string) *UserAgent {
	if userAgent == "" {
		return nil
	}

	ua := useragent.Parse(userAgent)

	return &UserAgent{
		Device:  ua.Device,
		OS:      ua.OS,
		Browser: ua.Name,
	}
}
