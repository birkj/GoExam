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
	"strconv"
	"time"
)

type Client struct {
	pb.RouteClient
	ClientID         pb.ClientID
	LamportTimeStamp pb.LamportTimeStamp
	Logfile          os.File
}

//Sends heartbeat to server every 1 second
func SendHeartBeatToServer(client Client) {
	for {
		//New context is made each time, if it takes more than a second to get answer from server it aborts
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		ack, err := client.SendHeartBeats(ctx, &pb.HeartBeat{Time: &client.LamportTimeStamp, Clientid: &client.ClientID})
		if err != nil {
			log.Fatalf("No response from server: %v", err)
		}
		log.Println(ack.Status)
		WriteToLogFile(client, ack.Status)
		time.Sleep(time.Second * 1)
	}
}

func ConnectToServer(ctx context.Context, client Client) {
	ack, err := client.Connect(ctx, &pb.ConnectRequest{Id: client.ClientID.Id})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	WriteToLogFile(client, ack.Status)
	log.Println(ack.Status)
}

func CreateLogFile(id int64) os.File {
	file, error := os.Create("client" + strconv.FormatInt(id, 10) + "_log.txt")
	if error != nil {
		log.Fatalf("Failed to create file: %v", error)
	}
	return *file
}

func WriteToLogFile(client Client, txt string) {
	t := time.Now()
	d := time.Second
	client.Logfile.WriteString(t.Round(d).String() + ": " + txt + "\n")
}

func main() {

	//Initital setup
	//Get client ID
	id := flag.Int64("id", 0, "current environment")
	port := flag.Int64("client_port", 8080, "current environment")
	flag.Parse()

	//Set up connection
	conn, err := grpc.Dial("localhost:"+strconv.FormatInt(*port, 10), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := Client{
		pb.NewRouteClient(conn),
		pb.ClientID{Id: *id},
		pb.LamportTimeStamp{Lamport: 0},
		CreateLogFile(*id),
	}

	connect_ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ConnectToServer(connect_ctx, client)

	//Start heartbeats as go func
	go SendHeartBeatToServer(client)

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

		reply, err := client.SayHello(ctx, &pb.RequestText{Body: string(text), Client: &client.ClientID})
		if err != nil {
			log.Fatalf("could not say hello: %v", err)
		}
		//Print reply
		log.Printf(reply.GetBody())
		WriteToLogFile(client, reply.GetBody())

	}

	//Send heartbeat in another thread

}
