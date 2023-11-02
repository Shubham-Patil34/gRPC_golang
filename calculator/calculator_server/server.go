package main

import (
	"context"
	"fmt"
	calculatorpb "grpc/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

// Define a server struct that implements the gRPC service
type CalculatorServiceServer struct {
	calculatorpb.UnsafeCalculatorServiceServer
}

// Implement a method of the gRPC service
func (s *CalculatorServiceServer) Calculator(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "Calculator function is invoked by: ", req)
	result := 0
	switch req.GetOperator() {
	case "+":
		result = int(req.GetOperand_1()) + int(req.GetOperand_2())
	case "-":
		result = int(req.GetOperand_1()) - int(req.GetOperand_2())
	case "*":
		result = int(req.GetOperand_1()) * int(req.GetOperand_2())
	case "%":
		result = int(req.GetOperand_1()) % int(req.GetOperand_2())
	}
	return &calculatorpb.CalculatorResponse{Result: int32(result)}, nil
}

// PrimeFactorization implements com_grpc_calculatorpb.CalculatorServiceServer.
func (s *CalculatorServiceServer) PrimeFactorization(req *calculatorpb.PrimeFactorizationRequest, stream calculatorpb.CalculatorService_PrimeFactorizationServer) error {
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "PrimeFactorization function is invoked by: ", req)
	// Get the number
	num := req.GetNumber()

	res := &calculatorpb.PrimeFactorizationResponse{}
	// While the number is even
	for num%2 == 0 {
		res.Factor = 2
		stream.Send(res)
		num /= 2
	}

	// While number is odd
	for i := 3; int32(i*i) <= num; i += 2 {
		for num%int32(i) == 0 {
			res.Factor = int32(i)
			stream.Send(res)
			num /= int32(i)
		}
	}

	// if the number is prime number greater than 2
	if num > 2 {
		res.Factor = num
		stream.Send(res)
	}
	return nil
}

// CalculateAverage implements com_grpc_calculatorpb.CalculatorServiceServer.
func (*CalculatorServiceServer) CalculateAverage(stream calculatorpb.CalculatorService_CalculateAverageServer) error {
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "CalculateAverage function is invoked by client stream:")
	sum := 0
	count := 0

	for {

		msg, err := stream.Recv()
		if err == io.EOF {
			// We have reached end of stream
			return stream.SendAndClose(&calculatorpb.CalculateAverageResponse{
				Average: float64(sum) / float64(count),
			})
		}

		if err != nil {
			log.Fatalf("Failed to read from CalculateAverageRequest stream: %v", err)
		}
		sum += int(msg.GetNumber())
		count++

	}
}

// FindMaximum implements com_grpc_calculatorpb.CalculatorServiceServer.
func (*CalculatorServiceServer) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "FindMaximum function is invoked by client stream:")
	max := 0

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// We have reached end of stream
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to read stream from client: %v", err)
		}
		if res.GetNumber() > int32(max) {
			max = int(res.GetNumber())
		}

		err = stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: int32(max),
		})
		if err != nil {
			log.Fatalf("Failed to send response to client: %v", err)
		}
	}
	
}

func main() {
	fmt.Println("Hello from Bhaskaracharya :)")

	// Create a gRPC server
	server := grpc.NewServer()

	// Register the gRPC service
	calculatorpb.RegisterCalculatorServiceServer(server, &CalculatorServiceServer{})

	// Create connection
	conn, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Start the gRPC server
	fmt.Println("Server started listening on: 50051")
	if err := server.Serve(conn); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
