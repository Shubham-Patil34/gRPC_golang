package main

import (
	"context"
	"fmt"
	calculatorpb "grpc/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I'm a calculator")

	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := calculatorpb.NewCalculatorServiceClient(conn)

	// Unary
	// doUnary(client)

	// Server Streaming
	doServerStreaming(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient){
		// Make a gRPC request
		req := &calculatorpb.CalculatorRequest{
			Operand_1: 23,
			Operand_2: 12,
			Operator: "-",
	   }
   
	   resp, err := client.Calculator(context.Background(), req)
	   if err != nil {
		   log.Fatalf("Failed to call Calculator %v", err)
	   }
   
	   // Process the gRPC response
	   fmt.Printf("Response from Calculator: %d\n", resp.Result)
}

func doServerStreaming(client calculatorpb.CalculatorServiceClient){
	var input int32
	fmt.Print("Enter a number: ")
    _, err := fmt.Scanf("%d", &input)
	if err != nil {
		log.Fatalf("Failed to read value from client %v", err)
	}

	// Make a gRPC request
	req := &calculatorpb.PrimeFactorizationRequest{
		Number: int32(input),
	}

   respStream, err := client.PrimeFactorization(context.Background(), req)
   if err != nil {
	   log.Fatalf("Failed to call Calculator %v", err)
   }

   // Process the gRPC response
   for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			// We have reached the end of stream
			break
		}
		if err != nil {
			log.Fatalf("Failed to read stream response %v", err)
		}
 
		fmt.Printf("Response from PrimeFactorization: %d\n", msg.GetFactor())
   }
}