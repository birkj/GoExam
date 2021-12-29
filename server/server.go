package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "program/route"
	"strconv"
)

var port = flag.String("port", "8080", "The docker port of the server")

type server struct {
	pb.UnimplementedRouteServer
	connectedClients []string
}

type argError struct {
	_type string
	_desc string
}

func (e *argError) Error() string {
	return e._type + ": " + e._desc
}

func (s *server) BroadcastMessage(ctx context.Context, in *pb.RequestText) (*pb.GenericText, error) {
	return &pb.GenericText{Body: in.Body}, nil
}

func (s *server) SayHello(ctx context.Context, inText *pb.RequestText) (*pb.ReplyText, error) {

	//Show text from client
	log.Println("Client " + strconv.FormatInt(inText.Client.GetId(), 10) + ": " + inText.GetBody())

	//Tell client that their message was recived
	return &pb.ReplyText{Body: inText.Body + " from server"}, nil

}

func (s *server) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.Acknowledgement, error) {

	//Show that a new client has connected on server
	log.Println("Client " + strconv.FormatInt(in.Id, 10) + ": has connected")

	//Add client to servers list of clients when connecting

	for _, client := range s.connectedClients {
		if client == strconv.FormatInt(in.Id, 10) {
			return nil, &argError{"Name doublication", "Name already exists on server"}
		}
	}

	s.connectedClients = append(s.connectedClients, strconv.FormatInt(in.Id, 10))
	log.Println(s.connectedClients)

	//Answer client
	return &pb.Acknowledgement{Status: "Successfully connected"}, nil
}

func main() {
	//Make connected client slice
	server := server{
		connectedClients: make([]string, 0),
	}

	//Start server
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRouteServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
