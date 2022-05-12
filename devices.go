package mailosaur

import (
	"strings"
	"time"
)

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

func (s *DevicesService) List() (*DeviceListResult, error) {
	result, err := s.client.HttpGet(&DeviceListResult{}, "api/devices")
	return result.(*DeviceListResult), err
}

func (s *DevicesService) Create(deviceCreateOptions DeviceCreateOptions) (*Device, error) {
	result, err := s.client.HttpPost(&Device{}, "api/devices", deviceCreateOptions)
	return result.(*Device), err
}

func (s *DevicesService) Otp(query string) (*OtpResult, error) {
	if strings.Contains(query, "-") {
		result, err := s.client.HttpGet(&OtpResult{}, "api/devices/"+query+"/otp")
		return result.(*OtpResult), err
	}

	result, err := s.client.HttpPost(&OtpResult{}, "api/devices/otp", &DeviceCreateOptions{SharedSecret: query})
	return result.(*OtpResult), err
}

func (s *DevicesService) Delete(id string) error {
	return s.client.HttpDelete("api/devices/" + id)
}
