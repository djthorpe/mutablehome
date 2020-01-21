package ecovacs

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Ecovacs struct {
	Country      string
	AccountId    string
	PasswordHash string
	Bus          gopi.Bus
}

type ecovacs struct {
	country, continent, lang, timezone string
	accountId, passwordHash            string
	deviceId, resourceId               string
	publicKey                          *rsa.PublicKey
	client                             *http.Client
	meta                               url.Values
	userId, accessToken                string
	devices                            []home.EvovacsDevice
	bus                                gopi.Bus

	base.Unit
}

type Response struct {
	Code      string `json:"code"`
	Message   string `json:"msg"`
	Timestamp uint   `json:"time"`
}

type AccessToken struct {
	Response
	Data struct {
		UserId      string `json:"uid"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		Country     string `json:"country"`
		AccessToken string `json:"accessToken"`
	} `json:"data"`
}

type AuthCode struct {
	Response
	Data struct {
		AuthCode   string `json:"authCode"`
		EcovacsUid string `json:"ecovacsUid"`
	} `json:"data"`
}

type UserResponse struct {
	Result       string `json:"result"`
	ErrorCode    uint   `json:"errno"`
	ErrorMessage string `json:"error"`
}

type LoginResponse struct {
	UserResponse
	UserId   string `json:"userId"`
	Resource string `json:"resource"`
	Token    string `json:"token"`
	Last     uint   `json:"last"`
}

type DevicesResponse struct {
	UserResponse
	Devices []*device `json:"devices"`
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MAIN_URL_FORMAT   = "https://eco-{{.country}}-api.ecovacs.com/v1/private/{{.country}}/{{.lang}}/{{.deviceId}}/{{.appCode}}/{{.appVersion}}/{{.channel}}/{{.deviceType}}{{.path}}"
	USER_URL_FORMAT   = "https://users-{{.continent}}.ecouser.net:8000/user.do"
	ECOVACS_REALM     = "ecouser.net"
	ECOVACS_XMPP_PORT = 5223
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Ecovacs) Name() string { return "mutablehome/ecovacs" }

func (config Ecovacs) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(ecovacs)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *ecovacs) Init(config Ecovacs) error {
	// Check bus parameter
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("Bus")
	} else {
		this.bus = config.Bus
	}

	// Check country, continent and language
	country := strings.ToUpper(config.Country)
	if continent, exists := COUNTRY_ALPHA2_TO_CONTINENT_CODE[country]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("ecovacs.country")
	} else {
		this.country = strings.ToLower(country)
		this.continent = strings.ToLower(continent)
		this.lang = "en"
	}
	// Fudge for Australia - they probably meant Austria (!?!)
	if this.country == "au" {
		this.continent = "eu"
	}

	// Set timezone
	this.timezone = "GMT"

	// Check account, password
	if accountId := strings.TrimSpace(config.AccountId); accountId == "" {
		return gopi.ErrBadParameter.WithPrefix("ecovacs.email")
	} else {
		this.accountId = accountId
	}
	if passwordHash := strings.TrimSpace(config.PasswordHash); passwordHash == "" {
		return gopi.ErrBadParameter.WithPrefix("ecovacs.password")
	} else {
		this.passwordHash = passwordHash
	}

	// Set deviceId and resouceId
	this.deviceId = MD5String(time.Now().String())
	this.resourceId = this.deviceId[0:8]

	// Load public key
	if key, err := DecodePublicKey(); err != nil {
		return err
	} else {
		this.publicKey = key
	}

	// Set HTTP client
	this.client = &http.Client{}

	// Set signing metadata
	this.meta = url.Values{
		"country":    []string{strings.ToLower(this.country)},
		"lang":       []string{strings.ToLower(this.lang)},
		"deviceId":   []string{this.deviceId},
		"appCode":    []string{"i_eco_e"},
		"appVersion": []string{"1.3.5"},
		"channel":    []string{"c_googleplay"},
		"deviceType": []string{"1"},
	}

	// Success
	return nil
}

func (this *ecovacs) Close() error {
	// Close devices
	err := gopi.NewCompoundError()
	for _, d := range this.devices {
		err.Add(d.(*device).Disconnect())
	}
	if err.ErrorOrSelf() != nil {
		return err.ErrorOrSelf()
	}

	// Release resources
	this.publicKey = nil
	this.client = nil
	this.meta = nil
	this.devices = nil
	this.bus = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.Ecovacs

func (this *ecovacs) Authenticate() error {
	var token AccessToken
	var authCode AuthCode

	if account, err := Encrypt(this.publicKey, this.accountId); err != nil {
		return err
	} else if password, err := Encrypt(this.publicKey, this.passwordHash); err != nil {
		return err
	} else if uri, err := this.mainURL("/user/login"); err != nil {
		return err
	} else if status, response, err := this.callMain(uri, url.Values{
		"account":  []string{account},
		"password": []string{password},
	}); err != nil {
		return err
	} else if status != http.StatusOK {
		return gopi.ErrUnexpectedResponse.WithPrefix(http.StatusText(status))
	} else if err := json.Unmarshal(response, &token); err != nil {
		return err
	} else if token.Code == "1005" {
		return home.ErrAuthenticationError
	} else if token.Code != "0000" {
		return gopi.ErrUnexpectedResponse.WithPrefix(token.Code)
	} else if uri, err := this.mainURL("/user/getAuthCode"); err != nil {
		return err
	} else if status, response, err := this.callMain(uri, url.Values{
		"uid":         []string{token.Data.UserId},
		"accessToken": []string{token.Data.AccessToken},
	}); err != nil {
		return err
	} else if status != http.StatusOK {
		return gopi.ErrUnexpectedResponse.WithPrefix(http.StatusText(status))
	} else if err := json.Unmarshal(response, &authCode); err != nil {
		return err
	} else if token.Code != "0000" {
		return gopi.ErrUnexpectedResponse.WithPrefix(token.Code)
	} else if response, err := this.callUserLogin(token, authCode); err != nil {
		return err
	} else {
		this.accessToken = response.Token
		if response.UserId != token.Data.UserId {
			this.userId = response.UserId
		} else {
			this.userId = token.Data.UserId
		}
	}

	// Authentication success
	return nil
}

func (this *ecovacs) Devices() ([]home.EvovacsDevice, error) {
	if this.userId == "" || this.accessToken == "" {
		return nil, gopi.ErrInternalAppError
	} else if this.devices == nil {
		this.devices = make([]home.EvovacsDevice, 0, 1)
	}

	if len(this.devices) == 0 {
		if devices, err := this.callDevices(); err != nil {
			return nil, err
		} else {
			for _, device := range devices {
				device.source = this
				this.devices = append(this.devices, device)
			}
		}
	}

	return this.devices, nil
}

// Connect to a device to start reading messages
func (this *ecovacs) Connect(d home.EvovacsDevice) error {
	for _, e := range this.devices {
		if d == e {
			this.Log.Debug("Connect:", d)
			return d.(*device).Connect()
		}
	}
	return gopi.ErrNotFound.WithPrefix("Connect")
}

// Disconnect from a device to stop updating
func (this *ecovacs) Disconnect(d home.EvovacsDevice) error {
	for _, e := range this.devices {
		if d == e {
			this.Log.Debug("Disconect:", d)
			return d.(*device).Disconnect()
		}
	}
	return gopi.ErrNotFound.WithPrefix("Disconnect")
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ecovacs) String() string {
	return "<" + this.Log.Name() +
		" country=" + strconv.Quote(this.country) +
		" continent=" + strconv.Quote(this.continent) +
		" lang=" + strconv.Quote(this.lang) +
		" timezone=" + strconv.Quote(this.timezone) +
		" accountId=" + strconv.Quote(this.accountId) +
		" passwordHash=" + strconv.Quote(this.passwordHash) +
		" devices=" + fmt.Sprint(this.Devices()) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *ecovacs) mainURL(path string) (*url.URL, error) {
	var buf bytes.Buffer
	data := map[string]string{
		"path": path,
	}
	for key := range this.meta {
		data[key] = this.meta.Get(key)
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

func (this *ecovacs) userURL() (*url.URL, error) {
	var buf bytes.Buffer
	data := map[string]string{
		"continent": this.continent,
	}
	if templ, err := template.New("USER_URL_FORMAT").Parse(USER_URL_FORMAT); err != nil {
		return nil, err
	} else if err := templ.Execute(&buf, data); err != nil {
		return nil, err
	} else if url, err := url.Parse(buf.String()); err != nil {
		return nil, err
	} else {
		return url, nil
	}
}

func (this *ecovacs) callDevices() ([]*device, error) {
	var devices DevicesResponse
	if data, err := this.callUser("GetDeviceList", map[string]interface{}{
		"userid": this.userId,
		"auth": map[string]string{
			"with":     "users",
			"userid":   this.userId,
			"realm":    ECOVACS_REALM,
			"token":    this.accessToken,
			"resource": this.resourceId,
		},
	}); err != nil {
		return nil, err
	} else if err := json.Unmarshal(data, &devices); err != nil {
		return nil, nil
	} else if devices.Result != "ok" && devices.ErrorMessage != "" {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(devices.ErrorMessage)
	} else if devices.Result != "ok" {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(devices.Result)
	} else {
		return devices.Devices, nil
	}
}

func (this *ecovacs) callUserLogin(token AccessToken, auth AuthCode) (*LoginResponse, error) {
	var login LoginResponse
	if data, err := this.callUser("loginByItToken", map[string]interface{}{
		"country":  strings.ToUpper(token.Data.Country),
		"resource": this.resourceId,
		"realm":    ECOVACS_REALM,
		"userId":   token.Data.UserId,
		"token":    auth.Data.AuthCode,
	}); err != nil {
		return nil, err
	} else if err := json.Unmarshal(data, &login); err != nil {
		return nil, err
	} else if login.Result != "ok" && login.ErrorMessage != "" {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(login.ErrorMessage)
	} else if login.Result != "ok" {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(login.Result)
	} else {
		return &login, nil
	}
}

func (this *ecovacs) callUser(path string, params map[string]interface{}) ([]byte, error) {
	params["todo"] = path
	if args, err := json.Marshal(params); err != nil {
		return nil, err
	} else if url, err := this.userURL(); err != nil {
		return nil, err
	} else if req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(args)); err != nil {
		return nil, err
	} else {
		req.Header.Add("Content-Type", "application/json")
		this.Log.Debug(req.URL, path)
		if response, err := this.client.Do(req); err != nil {
			return nil, err
		} else if response.StatusCode != http.StatusOK {
			return nil, gopi.ErrUnexpectedResponse.WithPrefix(http.StatusText(response.StatusCode))
		} else {
			defer response.Body.Close()
			if body, err := ioutil.ReadAll(response.Body); err != nil {
				return nil, err
			} else {
				return body, nil
			}
		}
	}
}

func (this *ecovacs) callMain(uri *url.URL, params url.Values) (int, []byte, error) {
	params.Add("requestId", MD5String(time.Now().String()))
	if req, err := http.NewRequest("GET", uri.String(), nil); err != nil {
		return http.StatusBadRequest, nil, err
	} else {
		req.URL.RawQuery = SortedQuery(this.signParams(params))
		this.Log.Debug(req.URL)
		if response, err := this.client.Do(req); err != nil {
			return http.StatusInternalServerError, nil, err
		} else {
			defer response.Body.Close()
			if body, err := ioutil.ReadAll(response.Body); err != nil {
				return http.StatusInternalServerError, nil, err
			} else {
				return response.StatusCode, body, nil
			}
		}
	}
}

func (this *ecovacs) signParams(params url.Values) url.Values {
	// Add values which we need to sign
	sign := url.Values{}
	sign.Set("authTimeZone", this.timezone)
	sign.Set("authTimespan", fmt.Sprint(time.Now().UnixNano()/1000000))
	for key := range params {
		sign.Set(key, params.Get(key))
	}
	for key := range this.meta {
		sign.Set(key, this.meta.Get(key))
	}

	// Create signature
	text := CLIENT_KEY
	for _, key := range SortedKeys(sign) {
		text += fmt.Sprintf("%s=%s", key, sign.Get(key))
	}
	text += SECRET

	// Add parameters to query
	params.Set("authAppkey", CLIENT_KEY)
	params.Set("authSign", MD5String(text))
	params.Set("authTimeZone", sign.Get("authTimeZone"))
	params.Set("authTimespan", sign.Get("authTimespan"))

	// Return parameters
	return params
}

func (this *ecovacs) deviceError(d *device, err error) {
	this.Log.Error(fmt.Errorf("%v: %w", d.Address(), err))

	// Disconnect and then reconnect when any device error occurs
	if err := this.Disconnect(d); err != nil {
		this.Log.Error(fmt.Errorf("%v: %w", d.Address(), err))
	}
	if err := this.Connect(d); err != nil {
		this.Log.Error(fmt.Errorf("%v: %w", d.Address(), err))
	}
}

func SortedQuery(params url.Values) string {
	sorted := url.Values{}
	for _, key := range SortedKeys(params) {
		sorted.Add(key, params.Get(key))
	}
	return sorted.Encode()
}

func SortedKeys(params url.Values) []string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
