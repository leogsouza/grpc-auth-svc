package main

import (
	"fmt"
	"log"
	"net"

	"github.com/leogsouza/grpc-auth-svc/pkg/pb"
	"github.com/leogsouza/grpc-auth-svc/pkg/config"
	"github.com/leogsouza/grpc-auth-svc/pkg/utils"
	"github.com/leogsouza/grpc-auth-svc/pkg/db"
	"github.com/leogsouza/grpc-auth-svc/pkg/services"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Auth Service is on", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve", err)
	}
}
