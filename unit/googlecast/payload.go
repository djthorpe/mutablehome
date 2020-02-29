/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

// Ref: https://github.com/vishen/go-chromecast/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Payload interface {
	WithId(id int) Payload
}

type PayloadHeader struct {
	Type      string `json:"type"`
	RequestId int    `json:"requestId,omitempty"`
}

type SetVolumeRequest struct {
	PayloadHeader
	Volume volume `json:"volume"`
}

type LaunchAppRequest struct {
	PayloadHeader
	AppId string `json:"appId"`
}

type ReceiverStatusResponse struct {
	PayloadHeader
	Status struct {
		Applications []application `json:"applications"`
		Volume       volume        `json:"volume"`
	} `json:"status"`
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

func (this *PayloadHeader) WithId(id int) Payload {
	this.RequestId = id
	return this
}

func (this *SetVolumeRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}

func (this *LaunchAppRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}
