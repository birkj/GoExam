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

type Replica struct {
	id         int32
	port       int32
	connection pb.RouteClient
}

type Client struct {
	id int32
}

type Server struct {
	pb.UnimplementedRouteServer
	Replicas         []Replica
	ConnectedClients []Client
	Logfile          os.File
	ServerID         int32
	LeaderID         int32
	LeaderPort       int32
	Port             int32
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
	WriteToLog(*s, msg)

	//Tell client that their message was recived
	return &pb.ReplyText{Body: inText.Body + " from Server"}, nil

}

func (s *Server) SendHeartBeats(ctx context.Context, in *pb.HeartBeat) (*pb.Acknowledgement, error) {
	msg := "Heartbeat from: " + fmt.Sprint(in.Id) + ". Lamport Timestamp: " + fmt.Sprint(in.Time.Lamport)
	WriteToLog(*s, msg)
	return &pb.Acknowledgement{Status: "Recieved heartbeat"}, nil
}

func (s *Server) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.Acknowledgement, error) {

	//Show that a new client has connected on Server
	msg := "Client " + fmt.Sprint(in.Id) + ": has connected"
	WriteToLog(*s, msg)

	//Add client to servers list of clients when connecting

	for _, client := range s.ConnectedClients {
		if client.id == in.Id {
			return nil, &argError{"Name doublication", "Name already exists on Server"}
		}
	}

	s.ConnectedClients = append(s.ConnectedClients, Client{id: in.Id})
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

func WriteToLog(s Server, txt string) {
	//To show in terminal once running
	log.Println(txt)
	//To write to file
	t := time.Now()
	d := time.Second
	s.Logfile.WriteString(t.Round(d).String() + ": " + txt + "\n")
}

func InitReplicas(s Server, numReplicas int32, port int32) {
	//Init clients from this server to all other servers, ports assumed to be fixed from program init
	WriteToLog(s, "Connecting to replicas.....")
	for i := 0; i < int(numReplicas); i++ {
		//Don't add yourself
		replicaPort := port + int32(i)
		if replicaPort == s.Port {
			continue
		}

		//Connect to replica port
		conn, err := grpc.Dial("localhost:"+fmt.Sprint(port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		s.Replicas = append(s.Replicas, Replica{id: int32(i), port: replicaPort, connection: pb.NewRouteClient(conn)})
		WriteToLog(s, "Connected to replica with id: "+fmt.Sprint(i))
	}
}

func main() {
	//Get arguments from flag
	id := flag.Int64("server_id", 0, "current environment")
	numReplicas := flag.Int64("num_replicas", 0, "current environment")
	leaderPort := flag.Int64("server_port_leader", 8080, "current environment")
	port := flag.Int64("server_port", 8080, "current environment")
	leader := flag.Int64("server_leader", 0, "current environment")
	flag.Parse()

	//Make connected client slice and log file
	server := Server{
		Replicas:         make([]Replica, 0),
		ConnectedClients: make([]Client, 0),
		Logfile:          CreateLogFile(int32(*id)),
		ServerID:         int32(*id),
		LeaderID:         int32(*leader),
		LeaderPort:       int32(*leaderPort),
		Port:             int32(*port),
	}

	//Write that we started server
	WriteToLog(server, "Started server "+fmt.Sprint(*id)+" on port: "+fmt.Sprint(*port))

	//Connect to replicas
	go InitReplicas(server, int32(*numReplicas), int32(*leaderPort))

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
