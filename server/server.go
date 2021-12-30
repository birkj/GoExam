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
	"sync"
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
	LamportTimestamp pb.LamportTimeStamp
	Mutex            sync.Mutex
}

type argError struct {
	_type string
	_desc string
}

func (e *argError) Error() string {
	return e._type + ": " + e._desc
}

//----------IMPLEMENTED FROM PROTO FILE

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

func (s *Server) SendLeaderHeartBeat(ctx context.Context, in *pb.HeartBeat) (*pb.Acknowledgement, error) {
	WriteToLog(*s, "Heartbeat from replica "+fmt.Sprint(in.Id)+" | Lamport: "+fmt.Sprint(in.Time.Lamport))
	return &pb.Acknowledgement{Status: "Leader (server " + fmt.Sprint(s.ServerID) + ") is alive"}, nil
}

func (s *Server) ElectionResult(ctx context.Context, in *pb.RequestText) (*pb.Acknowledgement, error) {
	//Delete old leader from replicas
	RemoveDeadLeader(s)

	//Set new leader
	WriteToLog(*s, "Update received: new leader is"+fmt.Sprint(in.Id))
	s.LeaderID = in.Id
	return &pb.Acknowledgement{Status: "Set new leader to " + fmt.Sprint(in.Id)}, nil
}

//-------------Local Functions------------
//Log file

func CreateLogFile(id int32) os.File {
	file, error := os.Create("server_" + fmt.Sprint(id) + "_log.txt")
	if error != nil {
		log.Fatalf("Failed to create file: %v", error)
	}
	return *file
}

//WriteToLog writes to an external file and logs to running terminal
func WriteToLog(s Server, txt string) {
	//To show in terminal once running
	log.Println(txt)
	//To write to file
	t := time.Now()
	d := time.Second
	s.Logfile.WriteString(t.Round(d).String() + ": " + txt + "\n")
}

//InitReplicas inits clients from this server to all other servers, ports assumed to be fixed from program init
func InitReplicas(s *Server, numReplicas int32, port int32) {
	WriteToLog(*s, "Connecting to replicas.....")
	for i := 0; i < int(numReplicas); i++ {
		//Don't add yourself
		replicaPort := port + int32(i)
		if replicaPort == s.Port {
			continue
		}

		//Connect to replica port
		conn, err := grpc.Dial("localhost:"+fmt.Sprint(replicaPort), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		s.Replicas = append(s.Replicas, Replica{id: int32(i), port: replicaPort, connection: pb.NewRouteClient(conn)})

		//Increase Lamport since this is an event
		s.LamportTimestamp.Lamport++
		WriteToLog(*s, "Connected to replica with id: "+fmt.Sprint(i)+" and port: "+fmt.Sprint(replicaPort)+" | Lamport: "+fmt.Sprint(s.LamportTimestamp.Lamport))
	}
}

//LeaderHeartBeat makes all replicas send a heartbeat to the current leader, the leader responds with an ack to confirm he/she is alive
func LeaderHeartBeat(s *Server) {
	for {
		//Leader doesn't ping itself
		if s.ServerID != s.LeaderID {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			for _, replica := range s.Replicas {
				//Replicas ping leader
				if replica.id == s.LeaderID {
					ack, err := replica.connection.SendLeaderHeartBeat(ctx, &pb.HeartBeat{Id: s.ServerID, Time: &s.LamportTimestamp})
					if err != nil {
						//ELECTION if no leader heartbeat
						WriteToLog(*s, "Leader not responding... calling election")
						CallElection(s)
						continue
					}
					WriteToLog(*s, ack.Status)
				}
			}
		}
		//Pings every 5 seconds (in order not to bloat the log, would be more often in a real system maybe every 150ms - 300ms or so)
		time.Sleep(time.Second * 5)
	}

}

//CallElection is called if the leader is not responding to
func CallElection(s *Server) {

	//Delete old leader from replicas
	RemoveDeadLeader(s)

	//Find biggest replica id
	max := int32(0)
	for _, replica := range s.Replicas {
		if replica.id > max {
			max = replica.id
		}
	}
	log.Println("MAX IS " + fmt.Sprint(max))
	//Leader is found, broadcast to others
	for _, replica := range s.Replicas {
		if replica.id == max {

			BroadcastElectionResult(s, replica)
		}
	}

}

func BroadcastElectionResult(s *Server, replica Replica) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, r := range s.Replicas {
		_, err := r.connection.ElectionResult(ctx, &pb.RequestText{Id: replica.id, Body: "New leader elected"})
		if err != nil {
			log.Println("Failed to send election result to server with id " + fmt.Sprint(replica.id) + " assuming it's dead.")
			//TODO: REMOVE DEAD REPLICA? OR ALLOW TO RECONNECT?
		}
	}
	s.LeaderID = replica.id
}

//RemoveReplica removes dead replicas (used to remove dead leader from Replica's list of other Replicas)
func RemoveReplica(s []Replica, index int) []Replica {
	return append(s[:index], s[index+1:]...)
}

func RemoveDeadLeader(s *Server) {
	//Remove dead server (old leader)
	for i, replica := range s.Replicas {
		if replica.id == s.LeaderID {
			s.Replicas = RemoveReplica(s.Replicas, i)
			break
		}
	}
}

//CalculateLamport finds max lamport and increase by 1
func CalculateLamport(s Server, other pb.LamportTimeStamp) *pb.LamportTimeStamp {

	new := pb.LamportTimeStamp{Lamport: 0}
	if s.LamportTimestamp.Lamport > other.Lamport {
		new.Lamport = s.LamportTimestamp.Lamport + 1
	} else {
		new.Lamport = other.Lamport + 1
	}
	return &new
}

//Since this is passive replication, a "server" has to act both as a client and a server to other "servers"
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
		LamportTimestamp: pb.LamportTimeStamp{Lamport: 0},
	}

	//Write that we started server
	WriteToLog(server, "Started server with id "+fmt.Sprint(*id)+" on port: "+fmt.Sprint(*port)+" | Lamport "+fmt.Sprint(server.LamportTimestamp.Lamport))

	//Connect to replicas
	go InitReplicas(&server, int32(*numReplicas), int32(*leaderPort))
	go LeaderHeartBeat(&server)

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
