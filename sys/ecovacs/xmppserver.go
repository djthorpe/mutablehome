package ecovacs

import (
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	ServerMap = map[string]string{
		"msg.ecouser.net":    "CH",
		"msg-as.ecouser.net": "TW,MY,JP,SG,TH,HK,IN,KR",
		"msg-na.ecouser.net": "US",
		"msg-eu.ecouser.net": "FR,ES,UK,NO,MX,DE,PT,CH,AU,IT,NL,SE,BE,DK",
		"msg-ww.ecouser.net": "",
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// CountryToXMPPServer returns hostname for countrycode
func CountryToXMPPServer(country string) string {
	d := ""
	for key, value := range ServerMap {
		if value == "" {
			d = key
		}
		for _, code := range strings.Split(value, ",") {
			if strings.ToUpper(country) == code {
				return key
			}
		}
	}
	return d
}
