/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	coap "github.com/go-ocf/go-coap"
	"github.com/go-ocf/go-coap/codes"
	dtls "github.com/pion/dtls/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPE CONVERSIONS

func boolToUint(value bool) uint {
	if value {
		return 1
	} else {
		return 0
	}
}

func durationToTransition(duration time.Duration) float32 {
	return float32(duration.Milliseconds()) / 100.0
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// coapConnectWith creates a secure COAP connection
func coapConnectWith(addr, key, value string, timeout time.Duration) (*coap.ClientConn, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: Missing key parameter", gopi.ErrBadParameter)
	} else if conn, err := coap.DialDTLSWithTimeout("udp", addr, &dtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			return []byte(value), nil
		},
		PSKIdentityHint: []byte(key),
		CipherSuites:    []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
	}, timeout); err != nil {
		return nil, fmt.Errorf("%w (addr: %s)", err, addr)
	} else {
		return conn, nil
	}
}

// coapAuthenticate performs the gateway authentication and returns JSON response
func coapAuthenticate(conn *coap.ClientConn, id string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Set identity
	body := strings.NewReader(fmt.Sprintf(`{"%s":"%s"}`, ATTR_IDENTITY, id))
	if response, err := conn.PostWithContext(ctx, PATH_AUTH_EXCHANGE, coap.AppJSON, body); err != nil {
		return nil, err
	} else if response.Code() != codes.Created {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(fmt.Sprint(response.Code()))
	} else {
		// Return success
		payload := response.Payload()
		return payload[0 : len(payload)-3], nil
	}
}

// addrForService returns an address to connect to from a service record
func addrForService(service gopi.RPCServiceRecord, flag gopi.RPCFlag) (string, error) {
	if service.Port == 0 {
		service.Port = DEFAULT_PORT
	}
	if flag == 0 {
		flag = gopi.RPC_FLAG_INET_V4 | gopi.RPC_FLAG_INET_V6
	}
	for _, addr := range service.Addrs {
		if addr == nil {
			continue
		} else if ip4 := addr.To4(); ip4 == nil && flag&gopi.RPC_FLAG_INET_V6 == gopi.RPC_FLAG_INET_V6 {
			return fmt.Sprintf("[%s]:%d", addr.String(), service.Port), nil
		} else if ip4 != nil && flag&gopi.RPC_FLAG_INET_V4 == gopi.RPC_FLAG_INET_V4 {
			return fmt.Sprintf("%s:%d", addr.String(), service.Port), nil
		}
	}

	// Return error
	return "", fmt.Errorf("%w: Cannot determine gateway address and/or port", gopi.ErrBadParameter)
}
