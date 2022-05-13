package mailosaur

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	apiKey := os.Getenv("MAILOSAUR_API_KEY")
	baseUrl := os.Getenv("MAILOSAUR_BASE_URL")

	if len(apiKey) == 0 {
		log.Fatal("Missing necessary environment variables - refer to README.md")
	}

	if len(baseUrl) == 0 {
		baseUrl = "https://next.mailosaur.com/"
	}

	client = New(apiKey)
	client.baseUrl = baseUrl
}

func TestDevicesCrud(t *testing.T) {
	deviceName := "My GO test"
	sharedSecret := "ONSWG4TFOQYTEMY="

	// Create a new device
	options := DeviceCreateOptions{
		Name:         deviceName,
		SharedSecret: sharedSecret,
	}

	createdDevice, err := client.Devices.Create(options)
	assert.NoError(t, err)

	assert.False(t, len(createdDevice.Id) == 0)
	assert.Equal(t, deviceName, createdDevice.Name)

	otpResult, err := client.Devices.Otp(createdDevice.Id)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(otpResult.Code))

	before, _ := client.Devices.List()
	assert.True(t, contains_device(before.Items, createdDevice.Id))

	client.Devices.Delete(createdDevice.Id)

	after, _ := client.Devices.List()
	assert.False(t, contains_device(after.Items, createdDevice.Id))
}

func TestOtpViaSharedSecret(t *testing.T) {
	sharedSecret := "ONSWG4TFOQYTEMY="

	otpResult, err := client.Devices.Otp(sharedSecret)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(otpResult.Code))
}

func contains_device(list []*Device, id string) bool {
	for _, x := range list {
		if id == x.Id {
			return true
		}
	}
	return false
}
