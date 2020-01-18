package ecovacs

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	xmpp "github.com/mattn/go-xmpp"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	DeviceId string `json:"did"`
	Name     string `json:"name"`
	Class    string `json:"class"`
	Resource string `json:"resource"`
	Nickname string `json:"nick"`
	Company  string `json:"company"`

	source *ecovacs
	client *xmpp.Client
}

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
// CONNECT

func (this *device) Connect() error {
	if server := CountryToXMPPServer(this.source.country); server == "" {
		return gopi.ErrInternalAppError.WithPrefix("country")
	} else {
		config := xmpp.Options{
			Host:     server,
			User:     fmt.Sprintf("%s@%s", this.source.userId, ECOVACS_REALM),
			Password: fmt.Sprintf("0/%s/%s", this.source.resourceId, this.source.accessToken),
			NoTLS:    true,
			Debug:    true,
			Session:  true,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		if client, err := config.NewClient(); err != nil {
			return err
		} else {
			this.client = client
		}
	}

	// Success
	return nil
}

func (this *device) Address() string {
	return fmt.Sprintf("%s@%s.ecorobot.net/atom", this.DeviceId, this.Class)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	return "<ecovacs.device" +
		" addr=" + strconv.Quote(this.Address()) +
		" nickname=" + strconv.Quote(this.Nickname) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

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
