package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"testing"
)

func TestConnection(t *testing.T) {

	kHost := os.Getenv("KEYSTONE_SERVICE_HOST")
	kPort := os.Getenv("KEYSTONE_SERVICE_PORT")
	if kHost == "" {
		kHost = "127.0.0.1"
	}
	if kPort == "" {
		kPort = "50031"
	}

	ksGrpcConn, err := grpc.Dial(kHost+":"+kPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Println(kHost + ":" + kPort)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	ksClient := proto.NewKeystoneClient(ksGrpcConn)
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	c.RegisterTypes(testSchemaType{})
	c.SyncSchema().Wait()

	log.Println("Marshalling")
	actor.Marshal(testSchemaType{})
}
