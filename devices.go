package mailosaur

import (
	"strings"
	"time"
)

// DevicesService provides operations for managing virtual security devices and retrieving
// their current one-time passwords (OTPs), used to automate testing of app-based
// multi-factor authentication. Accessed via the Devices field of MailosaurClient.
type DevicesService struct {
	client *MailosaurClient
}

type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DeviceListResult struct {
	Items []*Device `json:"items"`
}

type DeviceCreateOptions struct {
	Name         string `json:"name"`
	SharedSecret string `json:"sharedSecret"`
}

type OtpResult struct {
	Code    string    `json:"code"`
	Expires time.Time `json:"expires"`
}

// List returns a list of your virtual security devices. It returns a DeviceListResult
// containing your devices.
func (s *DevicesService) List() (*DeviceListResult, error) {
	result, err := s.client.HttpGet(&DeviceListResult{}, "api/devices")
	return result.(*DeviceListResult), err
}

// Create creates a new virtual security device. The deviceCreateOptions argument specifies
// the options used to create the device. It returns the newly-created Device.
func (s *DevicesService) Create(deviceCreateOptions DeviceCreateOptions) (*Device, error) {
	result, err := s.client.HttpPost(&Device{}, "api/devices", deviceCreateOptions)
	return result.(*Device), err
}

// Otp retrieves the current one-time password for a saved device, or for a given
// base32-encoded shared secret. The query argument is either the unique identifier of the
// device or a base32-encoded shared secret. It returns an OtpResult containing the current
// one-time password.
func (s *DevicesService) Otp(query string) (*OtpResult, error) {
	if strings.Contains(query, "-") {
		result, err := s.client.HttpGet(&OtpResult{}, "api/devices/"+query+"/otp")
		return result.(*OtpResult), err
	}

	result, err := s.client.HttpPost(&OtpResult{}, "api/devices/otp", &DeviceCreateOptions{SharedSecret: query})
	return result.(*OtpResult), err
}

// Delete permanently deletes a virtual security device. This operation cannot be undone. The
// id argument is the unique identifier of the device to delete.
func (s *DevicesService) Delete(id string) error {
	return s.client.HttpDelete("api/devices/" + id)
}
