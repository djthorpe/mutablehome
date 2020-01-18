package ecovacs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"gosrc.io/xmpp"
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
		config := xmpp.Config{
			TransportConfiguration: xmpp.TransportConfiguration{
				Address: fmt.Sprintf("%s:%v", server, ECOVACS_XMPP_PORT),
			},
			Jid:          fmt.Sprintf("%s@%s", this.source.userId, ECOVACS_REALM),
			Credential:   xmpp.Password(fmt.Sprintf("0/%s/%s", this.source.resourceId, this.source.accessToken)),
			StreamLogger: os.Stdout,
			Insecure:     true,
		}
		router := xmpp.NewRouter()
		if client, err := xmpp.NewClient(config, router, errorHandler); err != nil {
			return err
		} else {
			cm := xmpp.NewStreamManager(client, nil)
			log.Fatal(cm.Run())
		}

	}

	// Success
	return nil
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	return "<ecovacs.device" +
		" device_id=" + strconv.Quote(this.DeviceId) +
		" nickname=" + strconv.Quote(this.Nickname) +
		" company=" + strconv.Quote(this.Company) +
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
