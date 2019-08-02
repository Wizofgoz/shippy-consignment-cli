package main

import (
	"encoding/json"
	"errors"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/micro/cmd"
	pb "github.com/wizofgoz/shippy-consignment-service/proto/consignment"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
)

const (
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &consignment)

	return consignment, err
}

func main() {
	cmd.Init()

	// Create new greeter client
	client := pb.NewConsignmentServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	// Contact the server and print out its response.
	file := defaultFilename
	var token string
	log.Println(os.Args)

	if len(os.Args) < 3 {
		log.Fatal(errors.New("not enough arguments, expecting file and token"))
	}

	file = os.Args[1]
	token = os.Args[2]

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	r, err := client.Create(ctx, consignment)

	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}

	log.Printf("Created: %t", r.Created)

	getAll, err := client.Get(ctx, &pb.GetRequest{})

	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}

	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}