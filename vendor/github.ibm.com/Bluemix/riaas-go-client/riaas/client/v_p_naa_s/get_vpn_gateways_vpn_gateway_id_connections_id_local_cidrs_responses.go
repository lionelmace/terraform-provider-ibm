// Code generated by go-swagger; DO NOT EDIT.

package v_p_naa_s

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.ibm.com/Bluemix/riaas-go-client/riaas/models"
)

// GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsReader is a Reader for the GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrs structure.
type GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 404:
		result := NewGetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK creates a GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK with default headers values
func NewGetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK() *GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK {
	return &GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK{}
}

/*GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK handles this case with default header values.

The CIDRs were retrieved successfully.
*/
type GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK struct {
	Payload *models.VPNGatewayConnectionLocalCIDRs
}

func (o *GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK) Error() string {
	return fmt.Sprintf("[GET /vpn_gateways/{vpn_gateway_id}/connections/{id}/local_cidrs][%d] getVpnGatewaysVpnGatewayIdConnectionsIdLocalCidrsOK  %+v", 200, o.Payload)
}

func (o *GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.VPNGatewayConnectionLocalCIDRs)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound creates a GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound with default headers values
func NewGetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound() *GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound {
	return &GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound{}
}

/*GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound handles this case with default header values.

A resource with the specified identifier could not be found.
*/
type GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound struct {
	Payload *models.Riaaserror
}

func (o *GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound) Error() string {
	return fmt.Sprintf("[GET /vpn_gateways/{vpn_gateway_id}/connections/{id}/local_cidrs][%d] getVpnGatewaysVpnGatewayIdConnectionsIdLocalCidrsNotFound  %+v", 404, o.Payload)
}

func (o *GetVpnGatewaysVpnGatewayIDConnectionsIDLocalCidrsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Riaaserror)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
