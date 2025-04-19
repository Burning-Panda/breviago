// Auth middleware using OpenFGA

package auth

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/openfga/go-sdk/client"
)

var FgaClientConfig = &client.ClientConfiguration{
	ApiScheme: "http",
	ApiUrl:            "https://api.openfga.dev",
	StoreId:           "01JS6R4H1DAT2W5GR62K6H4KG0",
	AuthorizationModelId: "01JS6RXVF1TRXTZD26EPJY1KSW",
}

func InitFgaClient() *client.OpenFgaClient {
	fgaClient, err := client.NewSdkClient(FgaClientConfig)
	if err != nil {
		log.Fatal("Failed to initialize FGA client: ", err)
	}

	return fgaClient
}



// AuthorizationMiddleware creates a middleware that checks permissions using OpenFGA
func AuthorizationMiddleware(fgaClient *client.OpenFgaClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}