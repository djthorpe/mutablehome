package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	uuid "github.com/google/uuid"
)

///////////////////////////////////////////////////////////////////////////////

type Server struct {
	ssl       []string
	log       gopi.Logger
	debug     bool
	handlers  []*Handler
	client_id string

	http.Server
	GoogleActions
}

type Handler struct {
	method  string
	re      *regexp.Regexp
	handler HandlerFunc
}

type HandlerFunc func(resp http.ResponseWriter, req *http.Request)

///////////////////////////////////////////////////////////////////////////////

func NewServer(port uint, client_id string, logger gopi.Logger) *Server {
	// Check incoming parameters
	if port == 0 || logger == nil {
		return nil
	}

	// Create object
	this := new(Server)
	this.log = logger
	this.debug = logger.IsDebug()
	this.handlers = make([]*Handler, 0, 10)
	this.client_id = client_id
	this.agent_user_id = "1836.15267389"
	this.Addr = fmt.Sprintf(":%v", port)
	this.Handler = this

	// Create handler
	this.registerHandler(http.MethodGet, regexp.MustCompile("^/oauth2$"), this.AuthorizationHandler)
	this.registerHandler(http.MethodPost, regexp.MustCompile("^/mutablehome$"), this.GoogleActionHandler)

	// Success
	return this
}

func (this *Server) Serve() error {
	if len(this.ssl) == 0 {
		return this.Server.ListenAndServe()
	} else if len(this.ssl) == 2 {
		return this.Server.ListenAndServeTLS(this.ssl[0], this.ssl[1])
	} else {
		return gopi.ErrBadParameter
	}
}

func (this *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Debug request
	if this.debug {
		if dump, err := httputil.DumpRequest(req, true); err != nil {
			this.log.Error("ServeHTTP: %v", err)
		} else {
			fmt.Fprintln(os.Stderr, string(dump))
		}
	}
	// Handle request
	for _, handler := range this.handlers {
		if handler.method == req.Method {
			if matches := handler.re.FindStringSubmatch(req.URL.Path); len(matches) > 0 {
				handler.handler(resp, req)
				return
			}
		}
	}
	// Fallback to not found
	http.NotFound(resp, req)
}

///////////////////////////////////////////////////////////////////////////////

func (this *Server) registerHandler(method string, path *regexp.Regexp, handler HandlerFunc) error {
	this.handlers = append(this.handlers, &Handler{
		method, path, handler,
	})
	// Success
	return nil
}

func (this *Server) generateAccessToken(user_id string, create_at time.Time) string {
	buf := bytes.NewBufferString(this.client_id)
	buf.WriteString(user_id)
	buf.WriteString(strconv.FormatInt(create_at.UnixNano(), 10))
	access := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()
	access = strings.ToUpper(strings.TrimRight(access, "="))
	return access
}

///////////////////////////////////////////////////////////////////////////////

func (this *Server) AuthorizationHandler(w http.ResponseWriter, req *http.Request) {
	// Get all the parameters
	values := req.URL.Query()

	// We only allow the response_type=token (implicit) version here
	if values.Get("response_type") != "token" {
		this.log.Warn("Invalid response_type: %v", values.Get("reponse_type"))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Check client_id parameter
	if values.Get("client_id") != this.client_id {
		this.log.Warn("Invalid client_id: %v", values.Get("client_id"))
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}
	// Redirect URL
	if uri := values.Get("redirect_uri"); uri == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else if uri_, err := url.Parse(uri); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
	} else {
		uri_.Fragment = url.Values{
			"access_token": []string{this.generateAccessToken("1234", time.Now())},
			"token_type":   []string{"Bearer"},
			"expires_in":   []string{"3600"},
			"state":        values["state"],
		}.Encode()
		this.log.Debug("Redirect: %v", uri_)
		http.Redirect(w, req, uri_.String(), http.StatusTemporaryRedirect)
	}
}

///////////////////////////////////////////////////////////////////////////////

func (this *Server) GoogleActionHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: Check authorization
	if header := req.Header.Get("Authorization"); header == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else {
		this.log.Debug("TODO: Check: %v", header)
	}

	// Obtain body, handle action
	if err := this.GoogleActions.Handle(w, req); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
	}
}
