package main

import (
	"google.golang.org/grpc"
	"homework/internal/adapters/grpcapp"
	"homework/internal/adapters/osfilesys"
	"homework/internal/app"
	"homework/internal/ports/grpcserver"
	"log"
	"net"
	"os"
)

func main() {
	fileDir := os.Getenv("FILEDIR")
	port := os.Getenv("PORT")
	if fileDir == "" || port == "" {
		log.Fatal("file directory is not specified")
	}

	repo, err := osfilesys.NewFileRepository(fileDir)
	if err != nil {
		log.Fatal("unable to open server file directory")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpcserver.LoggingServerInterceptor))

	grpcapp.RegisterFileServiceServer(s, grpcserver.NewFileServiceServer(app.NewUseCase(repo)))
	log.Printf("server is listening at #%s port\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
