/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package tivo

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi-rpc"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type TiVo struct {
	MediaAccessKey string
}

type tivo struct {
	log gopi.Logger
	mak string
}

type session struct {
	conn           gopi.RPCClientConn
	schema_version uint32
	rpc_id         uint32
	session_id     uint32
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config TiVo) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("<tivo.Open>{ mak=%v }", strconv.Quote(config.MediaAccessKey))

	this := new(tivo)
	this.log = logger
	if mak := strings.TrimSpace(config.MediaAccessKey); mak == "" {
		return nil, fmt.Errorf("Missing -tivo.mak (Media Access Key) value")
	} else {
		this.mak = mak
	}

	// Success
	return this, nil
}

func (this *tivo) Close() error {
	this.log.Debug("<tivo.Close>{ mak=%v }", strconv.Quote(this.mak))

	// Return success
	return nil
}

func (this *tivo) NewSession(conn gopi.RPCClientConn) (rpc.TiVoSession, error) {
	session := new(session)
	session.schema_version = 17
	session.rpc_id = 0
	session.session_id = rand.Uint32()
	return session, nil
}

func (this *session) payload(fh io.Writer) error {

}

func (this *session) headers(fh io.Writer, req_type string, body_id string, response_count string) error {
	var header http.Header
	header.Add("Type", "request")
	header.Add("RpcId", fmt.Sprint(this.rpc_id))
	header.Add("SchemaVersion", fmt.Sprint(this.schema_version))
	header.Add("Content-Type", "application/json")
	header.Add("RequestType", req_type)
	header.Add("ResponseCount", response_count)
	header.Add("BodyId", body_id)
	header.Add("X-ApplicationSessionId", fmt.Sprint(this.session_id))
	return header.Write(fh)
}

func (this *session) preamble(headers, payload []byte) string {
	return fmt.Sprintf("MRPC/2 %v %v", len(headers)+2, len(payload))
}
