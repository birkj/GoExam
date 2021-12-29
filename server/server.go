package main

import (
	"context"
	"flag"
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
	msg := "Client " + strconv.FormatInt(inText.Client.GetId(), 10) + ": " + inText.GetBody()
	WriteToLogFile(*s, msg)
	log.Println(msg)

	//Tell client that their message was recived
	return &pb.ReplyText{Body: inText.Body + " from Server"}, nil

}

func (s *Server) SendHeartBeats(ctx context.Context, in *pb.HeartBeat) (*pb.Acknowledgement, error) {
	msg := "Heartbeat from: " + strconv.FormatInt(in.Clientid.Id, 10) + ". Lamport Timestamp: " + strconv.FormatInt(in.Time.Lamport, 10)
	WriteToLogFile(*s, msg)
	log.Println(msg)
	return &pb.Acknowledgement{Status: "Recieved heartbeat"}, nil
}

func (s *Server) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.Acknowledgement, error) {

	//Show that a new client has connected on Server
	msg := "Client " + strconv.FormatInt(in.Id, 10) + ": has connected"
	WriteToLogFile(*s, msg)
	log.Println(msg)

	//Add client to servers list of clients when connecting

	for _, client := range s.ConnectedClients {
		if client == strconv.FormatInt(in.Id, 10) {
			return nil, &argError{"Name doublication", "Name already exists on Server"}
		}
	}

	s.ConnectedClients = append(s.ConnectedClients, strconv.FormatInt(in.Id, 10))
	log.Println(s.ConnectedClients)

	//Answer client
	return &pb.Acknowledgement{Status: "Successfully connected"}, nil
}

func CreateLogFile(id int64) os.File {
	file, error := os.Create("server_" + strconv.FormatInt(id, 10) + "_log.txt")
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
		Logfile:          CreateLogFile(*id),
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
