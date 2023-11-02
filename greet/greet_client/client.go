package main

import (
	"context"
	"fmt"
	greetpb "grpc/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I'm a greeter")

	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := greetpb.NewGreetServiceClient(conn)

	// doUnary(client)
	// doServerStreaming(client)
	// doClientStreaming(client)
	doBiDiStreaming(client)
}

// Unary API function
func doUnary(client greetpb.GreetServiceClient){
	// Make a gRPC request
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FName: "Maheshrao",
			LName: "Jadhav",
		},
	}

	resp, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Greet: %v", err)
	}

	// Process the gRPC response
	fmt.Printf("Response from Greet: %s\n", resp.Result)
}

// Server Stream API function
func doServerStreaming(client greetpb.GreetServiceClient){
	// Make a gRPC request
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FName: "Maheshrao",
			LName: "Jadhav",
		},
	}

	respStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call GreetManyTimes: %v", err)
	}

	// Process the gRPC response
	for{
		msg, err := respStream.Recv()
		if err == io.EOF {
			// We have reached the end of stream
			break
		}
		if err != nil {
			log.Fatalf("Failed to read stream response from GreetManyTimes: %v", err)
		}
		log.Println("Response from GreetManyTimes:", msg.GetResult())
	}

}

// Server Stream API function
func doClientStreaming(client greetpb.GreetServiceClient){
	// Make a gRPC request
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FName: "Krishna",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FName: "Suraj",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FName: "Pruthvi",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FName: "Rohit",
			},
		},
	}

	respStream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Failed to call LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		respStream.Send(req)
		time.Sleep(2000 * time.Millisecond)
	}

	res, err := respStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to read response from LongGreet: %v", err)
	}
	fmt.Println("LongGreet response: ", res)
		
}

// Server Stream API function
func doBiDiStreaming(client greetpb.GreetServiceClient){
	// Make a gRPC request
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FName: "Krishna",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FName: "Suraj",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FName: "Pruthvi",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FName: "Rohit",
			},
		},
	}

	// Create a stream
	respStream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Failed to call GreetEveryone: %v", err)
	}

	waitc := make(chan struct{})
	// Send a bunch of messages (go routine)
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending Greeting: %v\n", req)
			respStream.Send(req)
			time.Sleep(2000 * time.Millisecond)
		}
		respStream.CloseSend()
	}()	
	
	// Receive a bunch of messages (go routine)
	go func() {
		for {
			res, err := respStream.Recv()
			if err == io.EOF {
				// We have reached the end of stream
				break
			}
			if err != nil {
				log.Fatalf("Failed to read response from GreetEveryone: %v", err)
			}
			fmt.Println("GreetEveryone response: ", res.GetResult())
		}
		close(waitc)
	}()

	// Block until everything is done
	<-waitc
		
}