package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	pb "program/route"
	"strconv"
	"time"
)

var port = flag.String("port", "8080", "The docker port of the Server")

type Server struct {
	pb.UnimplementedRouteServer
	ConnectedClients []string
	Logfile          os.File
	ServerID         int32
}

type argError struct {
	_type string
	_desc string
}

func (e *argError) Error() string {
	return e._type + ": " + e._desc
}

func (s *Server) BroadcastMessage(ctx context.Context, in *pb.RequestText) (*pb.GenericText, error) {
	return &pb.GenericText{Body: in.Body}, nil
}

func (s *Server) SayHello(ctx context.Context, inText *pb.RequestText) (*pb.ReplyText, error) {

	//Show text from client
	msg := "Client " + fmt.Sprint(inText.Client) + ": " + inText.GetBody()
	WriteToLogFile(*s, msg)
	log.Println(msg)

	//Tell client that their message was recived
	return &pb.ReplyText{Body: inText.Body + " from Server"}, nil

}

func (s *Server) SendHeartBeats(ctx context.Context, in *pb.HeartBeat) (*pb.Acknowledgement, error) {
	msg := "Heartbeat from: " + fmt.Sprint(in.Id) + ". Lamport Timestamp: " + fmt.Sprint(in.Time.Lamport)
	WriteToLogFile(*s, msg)
	log.Println(msg)
	return &pb.Acknowledgement{Status: "Recieved heartbeat"}, nil
}

func (s *Server) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.Acknowledgement, error) {

	//Show that a new client has connected on Server
	msg := "Client " + fmt.Sprint(in.Id) + ": has connected"
	WriteToLogFile(*s, msg)
	log.Println(msg)

	//Add client to servers list of clients when connecting

	for _, client := range s.ConnectedClients {
		if client == fmt.Sprint(in.Id) {
			return nil, &argError{"Name doublication", "Name already exists on Server"}
		}
	}

	s.ConnectedClients = append(s.ConnectedClients, fmt.Sprint(in.Id))
	log.Println(s.ConnectedClients)

	//Answer client
	return &pb.Acknowledgement{Status: "Successfully connected"}, nil
}

func CreateLogFile(id int32) os.File {
	file, error := os.Create("server_" + fmt.Sprint(id) + "_log.txt")
	if error != nil {
		log.Fatalf("Failed to create file: %v", error)
	}
	return *file
}

func WriteToLogFile(server Server, txt string) {
	t := time.Now()
	d := time.Second
	server.Logfile.WriteString(t.Round(d).String() + ": " + txt + "\n")
}

func main() {
	//Get arguments from flag
	id := flag.Int64("server_id", 0, "current environment")
	port := flag.Int64("server_port", 8080, "current environment")
	flag.Parse()

	//Make connected client slice and log file
	server := Server{

		ConnectedClients: make([]string, 0),
		Logfile:          CreateLogFile(int32(*id)),
		ServerID:         int32(*id),
	}

	//Start Server
	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(*port, 10))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRouteServer(s, &server)
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
