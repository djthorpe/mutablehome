/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type application struct {
	AppId        string `json:"appId"`
	DisplayName  string `json:"displayName"`
	IsIdleScreen bool   `json:"isIdleScreen"`
	SessionId    string `json:"sessionId"`
	StatusText   string `json:"statusText"`
	TransportId  string `json:"transportId"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this application) ID() string {
	return this.AppId
}

func (this application) Name() string {
	return this.DisplayName
}

func (this application) Status() string {
	return this.StatusText
}

func (this application) Equals(other application) bool {
	if this.AppId != other.AppId {
		return false
	}
	if this.DisplayName != other.DisplayName {
		return false
	}
	if this.IsIdleScreen != other.IsIdleScreen {
		return false
	}
	if this.SessionId != other.SessionId {
		return false
	}
	if this.StatusText != other.StatusText {
		return false
	}
	if this.TransportId != other.TransportId {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this application) String() string {
	return fmt.Sprintf("<cast.Application>{ id=%v name=%v status=%v session=%v transport=%v idle_screen=%v }",
		strconv.Quote(this.AppId), strconv.Quote(this.DisplayName), strconv.Quote(this.StatusText), strconv.Quote(this.SessionId), strconv.Quote(this.TransportId), this.IsIdleScreen)
}
