package ecovacs

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	stop   chan struct{}
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
			Host:     fmt.Sprintf("%s:%d", server, ECOVACS_XMPP_PORT),
			User:     fmt.Sprintf("%s@%s", this.source.userId, ECOVACS_REALM),
			Password: fmt.Sprintf("0/%s/%s", this.source.resourceId, this.source.accessToken),
			NoTLS:    true,
			Debug:    true,
			Session:  true,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		fmt.Printf("%+v\n", config)
		if client, err := config.NewClient(); err != nil {
			return err
		} else {
			this.client = client
			this.stop = make(chan struct{})
			go this.ping(this.stop)
			go func() {
				for {
					if stanza, err := this.client.Recv(); err != nil {
						fmt.Println("ERROR", err)
					} else {
						switch v := stanza.(type) {
						case xmpp.Chat:
							fmt.Println("CHAT", v.Remote, v.Text)
						case xmpp.Presence:
							fmt.Println("PRESENCE", v.From, v.Show)
						}
					}
				}
			}()
		}
	}

	// Success
	return nil
}

func (this *device) Disconnect() error {
	this.stop <- struct{}{}
	<-this.stop
	return nil
}

func (this *device) Address() string {
	return fmt.Sprintf("%s@%s.ecorobot.net/atom", this.DeviceId, this.Class)
}

func (this *device) FetchBatteryLevel() error {
	if _, err := this.send("1", `<ctl td="Clean"><clean type="auto" speed="standard"/></ctl>`); err != nil {
		return err
	} else {
		return nil
	}
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

func (this *device) ping(stop chan struct{}) {
	ticker := time.NewTimer(time.Second)
FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			fmt.Println("PING")
			if err := this.client.PingC2S(this.Address(), this.client.JID()); err != nil {
				fmt.Println("ERROR", err)
			}
			ticker.Reset(30 * time.Second)
		case <-stop:
			fmt.Println("STOP")
			break FOR_LOOP
		}
	}
	close(stop)
}

func (this *device) send(reqid, command string) (string, error) {
	return this.client.RawInformationQuery(this.Address(), this.client.JID(), reqid, xmpp.IQTypeSet, "com:ctl", command)
}
