package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	pb "program/route"
	"time"
)

func main() {

	//Get client ID
	id := flag.Int64("id", 0, "current environment")
	flag.Parse()

	//Set up connection
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewRouteClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//Connect to server
	ack, err := client.Connect(ctx, &pb.ConnectRequest{Id: *id})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(ack.Status)

	//Ask forever
	for {
		//Get text from input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _, _ := reader.ReadLine()

		//Send request and wait for reply
		//New context for every loop
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		reply, err := client.SayHello(ctx, &pb.RequestText{Body: string(text), Client: &pb.Client{Id: *id}})
		if err != nil {
			log.Fatalf("could not say hello: %v", err)
		}

		//Print reply
		log.Printf(reply.GetBody())
	}

}
