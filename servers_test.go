package mailosaur

import (
	"testing"
    "log"
    "os"
    assert "github.com/stretchr/testify/require"
)

func init() {
    apiKey := os.Getenv("MAILOSAUR_API_KEY")
    baseUrl := os.Getenv("MAILOSAUR_BASE_URL")

    if (len(apiKey) == 0) {
        log.Fatal("Missing necessary environment variables - refer to README.md")
    }

    if (len(baseUrl) == 0) {
        baseUrl = "https://next.mailosaur.com/"
    }

    client = New(apiKey)
	client.baseUrl = baseUrl
}

func TestList(t *testing.T) {
    result, err := client.Servers.List()
    assert.NoError(t, err)

    assert.True(t, len(result.Items) > 1)
}

func TestGetNotFound(t *testing.T) {
    // Should throw if server is not found
    _, err := client.Servers.Get("efe907e9-74ed-4113-a3e0-a3d41d914765")
    
    // TODO Assert is a MailosaurException
    assert.Error(t, err)
}

func TestCrud(t *testing.T) {
    serverName := "My GO test"

    // Create a new server
    options := ServerCreateOptions {
        Name: serverName,
    }

    createdServer, err := client.Servers.Create(options)    
    assert.NoError(t, err)

    assert.False(t, len(createdServer.Id) == 0)
    assert.Equal(t, serverName, createdServer.Name)
    assert.False(t, len(createdServer.Password) == 0)
    // assert.NotNil(t, createdServer.Users == 0)
    assert.Equal(t, 0, createdServer.Messages)

    retrievedServer, err := client.Servers.Get(createdServer.Id)
    assert.NoError(t, err)

    assert.Equal(t, createdServer.Id, retrievedServer.Id)
    assert.Equal(t, createdServer.Name, retrievedServer.Name)
    assert.False(t, len(retrievedServer.Password) == 0)
    // Assert.NotNull(retrievedServer.Users)
    assert.Equal(t, 0, retrievedServer.Messages)

    retrievedServer.Name += " updated with ellipsis … and emoji 👨🏿‍🚒";
    updatedServer, err := client.Servers.Update(retrievedServer.Id, retrievedServer)
    assert.NoError(t, err)

    assert.Equal(t, retrievedServer.Id, updatedServer.Id)
    assert.Equal(t, retrievedServer.Name, updatedServer.Name)
    assert.Equal(t, retrievedServer.Password, updatedServer.Password)
    // Assert.Equal(retrievedServer.Users, updatedServer.Users)
    assert.Equal(t, retrievedServer.Messages, updatedServer.Messages)

    err = client.Servers.Delete(retrievedServer.Id)
    assert.NoError(t, err)

    // Attempting to delete again should fail
    err = client.Servers.Delete(retrievedServer.Id)
    assert.Error(t, err)
    // TODO Assert is a MailosaurException
}

func TestFailedCreate(t *testing.T) {
    serverCreateOptions := ServerCreateOptions{}
    
    _, err := client.Servers.Create(serverCreateOptions)
    
    // TODO Assert is a MailosaurException
    assert.Error(t, err)

    // TODO Implement MailosaurException structure
    // Assert.Equal("Request had one or more invalid parameters.", ex.Message);
    // Assert.Equal("invalid_request", ex.ErrorType);
    // Assert.Equal(400, ex.HttpStatusCode);
    // Assert.Equal("{\"type\":\"ValidationError\",\"messages\":{\"name\":\"Please provide a name for your server\"}}", ex.HttpResponseBody);
}
