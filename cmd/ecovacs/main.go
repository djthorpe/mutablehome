package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"
)

const (
	// https://github.com/wpietri/sucks/blob/master/sucks/__init__.py
	CLIENT_KEY = "eJUWrzRv34qFSaYk"
	SECRET     = "Cyu5jcR4zyK6QEPn1hdIGXB5QIDAQABMA0GC"
	PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDb8V0OYUGP3Fs63E1gJzJh+7iq
eymjFUKJUqSD60nhWReZ+Fg3tZvKKqgNcgl7EGXp1yNifJKUNC/SedFG1IJRh5hB
eDMGq0m0RQYDpf9l0umqYURpJ5fmfvH/gjfHe3Eg/NTLm7QEa0a0Il2t3Cyu5jcR
4zyK6QEPn1hdIGXB5QIDAQAB
-----END PUBLIC KEY-----`
	MAIN_URL_FORMAT = "https://eco-{{.country}}-api.ecovacs.com/v1/private/{{.country}}/{{.lang}}/{{.deviceId}}/{{.appCode}}/{{.appVersion}}/{{.channel}}/{{.deviceType}}{{.path}}"
	USER_URL_FORMAT = "https://users-{continent}.ecouser.net:8000/user.do"
	REALM           = "ecouser.net"
)

var (
	ServerMap = map[string]string{
		"msg.ecouser.net":    "CH",
		"msg-as.ecouser.net": "TW,MY,JP,SG,TH,HK,IN,KR",
		"msg-na.ecouser.net": "US",
		"msg-eu.ecouser.net": "FR,ES,UK,NO,MX,DE,PT,CH,AU,IT,NL,SE,BE,DK",
		"msg-ww.ecouser.net": "",
	}
)

type ecovacs struct {
	Country      string
	Continent    string
	Lang         string
	AccountId    string
	DeviceId     string
	PasswordHash string

	publicKey *rsa.PublicKey
}

func NewEcovacs(country, email, passwordHash string) (*ecovacs, error) {
	this := new(ecovacs)
	if continent, exists := COUNTRY_ALPHA2_TO_CONTINENT_CODE[country]; exists == false {
		return nil, fmt.Errorf("Invalid -country")
	} else {
		this.Country = country
		this.Continent = continent
		this.Lang = "en"
	}
	// Fudge for AU (I think they thought it was austria?)
	if this.Country == "AU" {
		this.Continent = "EU"
	}
	// Check account, password
	if email = strings.TrimSpace(email); email == "" {
		return nil, fmt.Errorf("Invalid -email")
	} else {
		this.AccountId = email
	}
	if passwordHash = strings.TrimSpace(passwordHash); passwordHash == "" {
		return nil, fmt.Errorf("Invalid -password")
	} else {
		this.PasswordHash = passwordHash
	}

	// Set DeviceId
	this.DeviceId = MD5String(time.Now().String())

	// Load public key
	if block, _ := pem.Decode([]byte(PUBLIC_KEY)); block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	} else if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil, err
	} else {
		this.publicKey = pub.(*rsa.PublicKey)
	}

	// Success
	return this, nil
}

/*
func (this *ecovacs) Authenticate() error {

}

/*
func CountryToXMPPServer(country string) (string, string) {
	d := ""
	for key, value := range ServerMap {
		if value == "" {
			d = key
		}
		for _, code := range strings.Split(value, ",") {
			if strings.ToUpper(country) == code {
				return strings.ToLower(country), key
			}
		}
	}
	return strings.ToLower(country), d
}

*/

func (this *ecovacs) mainURL(path string) (*url.URL, error) {
	var buf bytes.Buffer
	data := map[string]string{
		"country":    strings.ToLower(this.Country),
		"lang":       strings.ToLower(this.Lang),
		"deviceId":   this.DeviceId,
		"appCode":    "i_eco_e",
		"appVersion": "1.3.5",
		"channel":    "c_googleplay",
		"deviceType": "1",
		"path":       path,
	}
	if path != "" && strings.HasPrefix(path, "/") == false {
		return nil, fmt.Errorf("Invalid path")
	}
	if templ, err := template.New("MAIN_URL_FORMAT").Parse(MAIN_URL_FORMAT); err != nil {
		return nil, err
	} else if err := templ.Execute(&buf, data); err != nil {
		return nil, err
	} else if url, err := url.Parse(buf.String()); err != nil {
		return nil, err
	} else {
		return url, nil
	}
}

func (this *ecovacs) callMain(path string, params url.Values) (string, error) {
	if u, err := this.mainURL(path); err != nil {
		return "", err
	} else {
		params.Add("requestId", MD5String(time.Now().String()))
		sign(params)

		u.RawQuery = params.Encode()
		if response, err := http.Get(u.String()); err != nil {
			return "", err
		} else {
			defer response.Body.Close()
			if body, err := ioutil.ReadAll(response.Body); err != nil {
				return "", err
			} else {
				fmt.Println(string(body))
			}
		}
	}
	return "", nil
}

func (this *ecovacs) Authenticate() error {
	if account, err := this.Encrypt(this.AccountId); err != nil {
		return err
	} else if password, err := this.Encrypt(this.PasswordHash); err != nil {
		return err
	} else if url, err := this.callMain("/user/login", url.Values{
		"account":  []string{account},
		"password": []string{password},
	}); err != nil {
		return err
	} else {
		fmt.Println(url)
	}
	return nil
}

func (this *ecovacs) Encrypt(data string) (string, error) {
	if data, err := rsa.EncryptPKCS1v15(rand.Reader, this.publicKey, []byte(data)); err != nil {
		return "", fmt.Errorf("Encrypt: %w", err)
	} else {
		return string(data), nil
	}
}

func MD5String(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", md5.Sum(nil))
}

func Main(app gopi.App, args []string) error {
	country := strings.ToUpper(app.Flags().GetString("country", gopi.FLAG_NS_DEFAULT))
	email := app.Flags().GetString("email", gopi.FLAG_NS_DEFAULT)
	passwordHash := MD5String(app.Flags().GetString("password", gopi.FLAG_NS_DEFAULT))
	if ecovacs, err := NewEcovacs(country, email, passwordHash); err != nil {
		return err
	} else if err := ecovacs.Authenticate(); err != nil {
		return err
	} else {
		fmt.Println(ecovacs)
	}

	// Return success
	return nil
}

func main() {
	if app, err := app.NewCommandLineTool(Main, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		app.Flags().FlagString("country", "au", "Ecovacs Country Code")
		app.Flags().FlagString("email", "", "Ecovacs Account Email")
		app.Flags().FlagString("password", "", "Ecovacs Account Password")
		os.Exit(app.Run())
	}
}
