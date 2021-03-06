/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	pb "github.com/djthorpe/mutablehome/protobuf/castchannel"
	proto "github.com/golang/protobuf/proto"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type channel struct {
	C         chan interface{}
	messageId int
	sync.Mutex
}

type mediaUrl struct {
	url      string
	mimetype string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	CAST_DEFAULT_SENDER   = "sender-0"
	CAST_DEFAULT_RECEIVER = "receiver-0"
	CAST_NS_CONN          = "urn:x-cast:com.google.cast.tp.connection"
	CAST_NS_HEARTBEAT     = "urn:x-cast:com.google.cast.tp.heartbeat"
	CAST_NS_RECV          = "urn:x-cast:com.google.cast.receiver"
	CAST_NS_MEDIA         = "urn:x-cast:com.google.cast.media"
)

////////////////////////////////////////////////////////////////////////////////
// CONNECT AND DISCONNECT MESSAGES

func (this *channel) Connect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CONNECT"}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

func (this *channel) Disconnect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CLOSE"}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

func (this *channel) GetStatus() (int, []byte, error) {
	payload := &PayloadHeader{Type: "GET_STATUS"}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

func (this *channel) ConnectMedia(transportId string) (int, []byte, error) {
	payload := &PayloadHeader{Type: "CONNECT"}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

func (this *channel) DisconnectMedia(transportId string) (int, []byte, error) {
	payload := &PayloadHeader{Type: "CLOSE"}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

func (this *channel) GetMediaStatus(transportId string) (int, []byte, error) {
	payload := &PayloadHeader{Type: "GET_STATUS"}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

func (this *channel) SetVolume(v volume) (int, []byte, error) {
	payload := &SetVolumeRequest{PayloadHeader{Type: "SET_VOLUME"}, v}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

func (this *channel) SetMuted(muted bool) (int, []byte, error) {
	payload := &SetVolumeRequest{PayloadHeader{Type: "SET_VOLUME"}, volume{Muted_: muted}}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

func (this *channel) LaunchAppWithId(appId string) (int, []byte, error) {
	payload := &LaunchAppRequest{PayloadHeader{Type: "LAUNCH"}, appId}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

func (this *channel) PlayStop(state bool) (int, []byte, error) {
	payload := &PayloadHeader{}
	switch state {
	case true:
		payload.Type = "PLAY"
	case false:
		payload.Type = "STOP"
	}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

func (this *channel) PlayPause(state bool) (int, []byte, error) {
	payload := &PayloadHeader{}
	switch state {
	case true:
		payload.Type = "PLAY"
	case false:
		payload.Type = "PAUSE"
	}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

func (this *channel) LoadUrl(transportId string, media mediaUrl, autoplay bool) (int, []byte, error) {
	payload := &LoadMediaRequest{}
	payload.PayloadHeader = PayloadHeader{Type: "LOAD"}
	payload.Autoplay = autoplay
	payload.Media.ContentId = media.url
	payload.Media.ContentType = media.mimetype
	payload.Media.StreamType = "BUFFERED"
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

func (this *channel) LoadPlaylist(transportId string, media []mediaUrl, repeat bool) (int, []byte, error) {
	payload := &LoadQueueRequest{}
	payload.PayloadHeader = PayloadHeader{Type: "QUEUE_LOAD"}
	if repeat {
		payload.RepeatMode = "REPEAT_ALL"
	} else {
		payload.RepeatMode = "REPEAT_OFF"
	}
	payload.Items = make([]LoadQueueItem, len(media))
	for i, item := range media {
		payload.Items[i].Autoplay = true
		payload.Items[i].Media.ContentId = item.url
		payload.Items[i].Media.ContentType = item.mimetype
		payload.Items[i].Media.StreamType = "BUFFERED"
	}
	id := this.nextMessageId()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

////////////////////////////////////////////////////////////////////////////////
// SEND MESSAGES

func (this *channel) encode(source, dest, ns string, payload Payload) ([]byte, error) {
	json, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	payloadStr := string(json)
	message := &pb.CastMessage{
		ProtocolVersion: pb.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &source,
		DestinationId:   &dest,
		Namespace:       &ns,
		PayloadType:     pb.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadStr,
	}
	proto.SetDefaults(message)
	return proto.Marshal(message)
}

////////////////////////////////////////////////////////////////////////////////
// RECEIVE MESSAGES

func (this *channel) decode(data []byte) ([]byte, error) {
	message := &pb.CastMessage{}
	if err := proto.Unmarshal(data, message); err != nil {
		return nil, err
	}
	ns := message.GetNamespace()
	switch ns {
	case CAST_NS_RECV:
		return this.rcvMessageReceiver(message)
	case CAST_NS_HEARTBEAT:
		return this.rcvMessageHeartbeat(message)
	case CAST_NS_CONN:
		return this.rcvMessageConnection(message)
	case CAST_NS_MEDIA:
		return this.rcvMessageMedia(message)
	default:
		return nil, fmt.Errorf("Ignoring message with namespace %v", strconv.Quote(ns))
	}
}

func (this *channel) rcvMessageReceiver(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}

	switch header.Type {
	case "RECEIVER_STATUS":
		var receiverStatus ReceiverStatusResponse
		if err := json.Unmarshal([]byte(message.GetPayloadUtf8()), &receiverStatus); err != nil {
			return nil, fmt.Errorf("RECEIVER_STATUS: %w", err)
		}

		// Emit the volume
		this.C <- receiverStatus.Status.Volume

		// Emit the first application (doesn't support more than one)
		if len(receiverStatus.Status.Applications) == 0 {
			this.C <- application{}
		} else {
			this.C <- receiverStatus.Status.Applications[0]
		}
	case "INVALID_REQUEST":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(message.GetPayloadUtf8())
	case "LAUNCH_ERROR":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(message.GetPayloadUtf8())
	default:
		return nil, fmt.Errorf("Ignoring message %v in namespace %v", strconv.Quote(header.Type), strconv.Quote(message.GetNamespace()))
	}

	// Return success
	return nil, nil
}

func (this *channel) rcvMessageHeartbeat(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	case "PING":
		payload := &PayloadHeader{Type: "PONG", RequestId: -1}
		return this.encode(message.GetDestinationId(), message.GetSourceId(), message.GetNamespace(), payload)
	default:
		return nil, fmt.Errorf("Ignoring message %v in namespace %v", strconv.Quote(header.Type), strconv.Quote(message.GetNamespace()))
	}
}

func (this *channel) rcvMessageConnection(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	default:
		return nil, fmt.Errorf("Ignoring message %v in namespace %v", strconv.Quote(header.Type), strconv.Quote(message.GetNamespace()))
	}
	// Return success
	return nil, nil
}

func (this *channel) rcvMessageMedia(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	case "MEDIA_STATUS":
		var mediaStatus MediaStatusResponse
		if err := json.Unmarshal([]byte(message.GetPayloadUtf8()), &mediaStatus); err != nil {
			return nil, fmt.Errorf("MEDIA_STATUS: %w", err)
		}
		// Emit the media items
		this.C <- mediaStatus.Status
	case "LOAD_FAILED":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(message.GetPayloadUtf8())
	default:
		return nil, fmt.Errorf("Ignoring message %v in namespace %v", strconv.Quote(header.Type), strconv.Quote(message.GetNamespace()))
	}
	// Return success
	return nil, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *channel) nextMessageId() int {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Cycle messages from 1 to 99999
	this.messageId = (this.messageId + 1) % 100000
	return this.messageId
}
