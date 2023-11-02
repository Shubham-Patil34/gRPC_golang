package main

import (
	"context"
	"fmt"
	calculatorpb "grpc/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doServerStreaming(client)

	// Client Streaming
	// doClientStreaming(client)
	// doBiDiStreaming(client)
	doSquareRoot(client)
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

func doClientStreaming(client calculatorpb.CalculatorServiceClient){
	nums := []int32{10, 20, 33, 40, 50, 60}

	respStream, err := client.CalculateAverage(context.Background())
	if err != nil {
		log.Fatalf("Failed to call CalculateAverage: %v", err)
	}

	for _, num := range nums {
		fmt.Printf("Sending number: %v\n", num)
		respStream.Send(&calculatorpb.CalculateAverageRequest{
			Number: int32(num),
		})
		time.Sleep(2000 * time.Millisecond)
	}

	res, err := respStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to read response from CalculateAverage: %v", err)
	}
	fmt.Println("CalculateAverage response: ", res.GetAverage())
}

func doBiDiStreaming(client calculatorpb.CalculatorServiceClient){
	nums := []int32{10, 20, 5, 200, 50, 420}

	respStream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Failed to call FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for _, num := range nums {
			fmt.Printf("Sending number: %v\n", num)
			respStream.Send(&calculatorpb.FindMaximumRequest{
				Number: int32(num),
			})
			time.Sleep(2000 * time.Millisecond)
		}
		respStream.CloseSend()
	}()

	go func() {
		for {
			res, err := respStream.Recv()
			if err == io.EOF {
				// We have reached the end of stream
				break
			}
			if err != nil {
			log.Fatalf("Failed to read response from FindMaximum: %v", err)
			}
			fmt.Println("FindMaximum response: ", res.GetMaximum())
		}
		close(waitc)
	}()

	<-waitc
}

// Error code implementation for Square root calculation
func doSquareRoot(client calculatorpb.CalculatorServiceClient){
	// Correct call
	doErrorCall(client, 10)

	// Erroneous  call
	doErrorCall(client, -3)
}

func doErrorCall(client calculatorpb.CalculatorServiceClient, number float64) {
	resp, err := client.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: number,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// Actual err from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Failed to send request to SquareRoot: %v", err)
		}
	}

	fmt.Println("Response from SquareRoot: ", resp.GetRoot())
}