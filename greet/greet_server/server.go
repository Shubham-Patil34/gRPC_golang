package main

import (
	"context"
	"fmt"
	greetpb "grpc/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

// Define a server struct that implements the gRPC service
type GreetServiceServer struct {
	greetpb.UnsafeGreetServiceServer
}

// Implement a method of the gRPC service
func (s *GreetServiceServer) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "Greet function is invoked by: ", req)
	message := fmt.Sprintf("Ram Ram, %s!", req.GetGreeting().GetFName())
	return &greetpb.GreetResponse{Result: message}, nil
}

// Implement a method of the gRPC service
func (s *GreetServiceServer) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "GreetManyTimes function was invoked with: ", req)
	firstName := req.GetGreeting().GetFName()

	for i := 0; i < 10; i++ {
		result := "Namaskar " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

// LongGreet implements com_grpc_greetpb.GreetServiceServer.
func (*GreetServiceServer) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "LongGreet function was invoked")
	result := "ale..!"
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// We have reached end of stream
			// Send respond to client
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Failed to receive stream from LongGreet: %v", err)
		}
		result = msg.Greeting.GetFName() + ", " + result
	}
	
}

// GreetEveryone implements com_grpc_greetpb.GreetServiceServer.
func (*GreetServiceServer) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "GreetEveryone function was invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// We have reached the end of stream
			stream.Send(&greetpb.GreetEveryoneResponse{
				Result: "Enjoy..!",
			})
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to read client stream: %v", err)
		}

		fmt.Println("Client greeting from : ", req.Greeting.GetFName())
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Radhe Radhe, "+ req.Greeting.GetFName(),
		})
		if sendErr != nil {
			log.Fatalf("Failed to send response to client: %v", sendErr)
		}
	}

}

func main() {
	fmt.Println("Namaskar")

	// Create a gRPC server
	server := grpc.NewServer()

	// Register the gRPC service
	greetpb.RegisterGreetServiceServer(server, &GreetServiceServer{})

	// Create connection
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Start the gRPC server
	log.Println("Server started and listening on: 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
