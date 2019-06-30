package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

///////////////////////////////////////////////////////////////////////////////

type GoogleActionRequest struct {
	Inputs    []*GoogleActionInput `json:"inputs"`
	RequestId string               `json:"requestId"`
}

type GoogleActionInput struct {
	Intent string `json:"intent"`
}

type GoogleSyncResponse struct {
	RequestId string             `json:"requestId"`
	Payload   *GoogleSyncPayload `json:"payload"`
}

type GoogleSyncPayload struct {
	AgentUserId string                `json:"agentUserId"`
	Devices     []*GoogleActionDevice `json:"devices,omitempty"`
}

type GoogleActionDevice struct {
	Id              string                  `json:"id"`
	Type            string                  `json:"type"`
	Traits          []string                `json:"traits"`
	Name            *GoogleActionDeviceName `json:"name,omitempty"`
	WillReportState bool                    `json:"willReportState"`
	Room            string                  `json:"roomHint,omitempty"`
	DeviceInfo      *GoogleActionDeviceInfo `json:"deviceInfo,omitempty"`
}

type GoogleActionDeviceName struct {
	Name         string   `json:"name,omitempty"`
	DefaultNames []string `json:"defaultNames,omitempty"`
	Nicknames    []string `json:"nicknames,omitempty"`
}

type GoogleActionDeviceInfo struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	HwVersion    string `json:"hwVersion,omitempty"`
	SwVersion    string `json:"swVersion,omitempty"`
}

type GoogleActions struct {
	agent_user_id string
	devices       []*GoogleActionDevice
}

///////////////////////////////////////////////////////////////////////////////

const (
	GOOGLE_INTENT_SYNC       = "action.devices.SYNC"
	GOOGLE_INTENT_QUERY      = "action.devices.QUERY"
	GOOGLE_INTENT_EXECUTE    = "action.devices.EXECUTE"
	GOOGLE_INTENT_DISCONNECT = "action.devices.DISCONNECT"
)

const (
	GOOGLE_TYPE_COFFEE_MAKER = "action.devices.types.COFFEE_MAKER"
)

const (
	GOOGLE_TRAIT_ONOFF = "action.devices.traits.OnOff"
)

///////////////////////////////////////////////////////////////////////////////

func (this *GoogleActions) AddDevice(device *GoogleActionDevice) error {
	if device == nil {
		return gopi.ErrBadParameter
	}

	// Append device
	if this.devices == nil {
		this.devices = make([]*GoogleActionDevice, 0, 1)
	}
	this.devices = append(this.devices, device)

	// Success
	return nil
}

func (this *GoogleActions) Handle(w http.ResponseWriter, r *http.Request) error {
	var req *GoogleActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	for _, action := range req.Inputs {
		switch action.Intent {
		case GOOGLE_INTENT_SYNC:
			if response, err := this.Sync(req); err != nil {
				return err
			} else if err := json.NewEncoder(w).Encode(response); err != nil {
				return err
			} else if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
				return err
			}
		case GOOGLE_INTENT_EXECUTE:
			if response, err := this.Execute(req); err != nil {
				return err
			} else if err := json.NewEncoder(w).Encode(response); err != nil {
				return err
			} else if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unhandled intent: %v", strconv.Quote(action.Intent))
		}
	}

	// Success
	return nil
}

func (this *GoogleActions) Sync(r *GoogleActionRequest) (*GoogleSyncResponse, error) {
	response := &GoogleSyncResponse{
		RequestId: r.RequestId,
		Payload: &GoogleSyncPayload{
			AgentUserId: this.agent_user_id,
			Devices:     this.devices,
		},
	}
	return response, nil
}

func (this *GoogleActions) Execute(r *GoogleActionRequest) (*GoogleExecuteResponse, error) {
	response := &GoogleSyncResponse{
		RequestId: r.RequestId,
		Payload: &GoogleSyncPayload{
			AgentUserId: this.agent_user_id,
			Devices:     this.devices,
		},
	}
	return response, nil
}
