package main

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"

	privatepb "github.com/dane/protoc-gen-go-svc/example/proto/go/private"
	servicepb "github.com/dane/protoc-gen-go-svc/example/proto/go/service"
	servicev2 "github.com/dane/protoc-gen-go-svc/example/proto/go/service/v2"
	publicpb "github.com/dane/protoc-gen-go-svc/example/proto/go/v2"
	private "github.com/dane/protoc-gen-go-svc/example/service/private"
	v2 "github.com/dane/protoc-gen-go-svc/example/service/v2"
)

func TestExample(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	fatalIf(t, err)

	server := grpc.NewServer()
	defer server.Stop()

	impl := &private.Service{Store: make(map[string]*privatepb.Person)}
	option := v2.Converter{servicev2.NewConverter()}
	servicepb.RegisterServer(server, impl, option)

	go func(t *testing.T, server *grpc.Server, ln net.Listener) {
		err := server.Serve(ln)
		fatalIf(t, err)
	}(t, server, ln)

	conn, err := grpc.Dial(ln.Addr().String(), grpc.WithInsecure())
	fatalIf(t, err)
	client := publicpb.NewPeopleClient(conn)

	var personId string
	t.Run("create person in v2", func(t *testing.T) {
		resp, err := client.Create(context.Background(), &publicpb.CreateRequest{
			FullName:   "Dane Harrigan",
			Age:        35,
			Employment: publicpb.Person_FULL_TIME,
		})
		fatalIf(t, err)

		personId = resp.Person.Id
	})

	t.Logf("personId: %s", personId)

}

func fatalIf(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
