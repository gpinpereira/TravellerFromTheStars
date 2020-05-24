// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "../services"
	"../universe"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var main_solar_system *universe.SolarSystem

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer

	system *universe.SolarSystem
}

// SayHello implements helloworld.GreeterServer
func (s *server) GetBodyPosition(ctx context.Context, in *pb.BodyName) (*pb.BodyPos, error) {

	//log.Printf("Received: %v", in.GetName())
	allbodies := main_solar_system.GetBodies()

	if bodyfound, ok := allbodies[in.GetName()]; ok {
		x, y := bodyfound.GetPos()
		//log.Println(in.GetName())
		//log.Println(x)
		//log.Println(y)
		return &pb.BodyPos{X: x, Y: y}, nil
	}
	/*for _, body := range otherbodies {

		if body.Name == in.GetName() {
			x, y := body.GetPos()
			log.Println(in.GetName())
			log.Println(x)
			log.Println(y)
		}

	}*/
	//var sr string = "Hello " + in.GetName()
	return &pb.BodyPos{X: -1, Y: -1}, nil
}

func (s *server) SolarSystemPositions(ctx context.Context, in *pb.BodyName) (*pb.AllBodies, error) {
	log.Println("request SolarSystemPositions")
	allbodies := main_solar_system.GetBodies()
	//log.Println("allbodies")
	send_bodies := make(map[string]*pb.BodyPos)
	for name, body := range allbodies {

		x, y := body.GetPos()
		send_bodies[name] = &pb.BodyPos{X: x, Y: y}
		/*if body.Name == in.GetName() {
			x, y := body.GetPos()
			log.Println(in.GetName())
			log.Println(x)
			log.Println(y)
		}*/

	}
	return &pb.AllBodies{Bodies: send_bodies}, nil
}

func (s *server) RequestSolarSystemStats(ctx context.Context, in *pb.BodyName) (*pb.SolarSystemStats, error) {

	return &pb.SolarSystemStats{Age: main_solar_system.GetAge()}, nil
}

/*func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	var sr string = "Hello again " + in.GetName()
	return &pb.HelloReply{Message: &sr}, nil
}*/

func StartUniverse() {

	//system := universe.MakeSystem()
	system := universe.MakeSystemCSV()
	main_solar_system = system
	go universe.SimulateSystem(system)

	//time.Sleep(5 * time.Second)
	fmt.Println("started Universe")
	//fmt.Println((*main_solar_system).GetAge())

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{system: main_solar_system})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func main() {

	StartUniverse()

}
