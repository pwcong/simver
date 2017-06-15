package rpc

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"net"

	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	Init "github.com/pwcong/simver/init"
	pb "github.com/pwcong/simver/vertify"
	"google.golang.org/grpc/reflection"

	Redis "github.com/pwcong/simver/db/redis"
	IndexService "github.com/pwcong/simver/service/index"
)

type Server struct {
	IP   string
	Port int
}

func (s *Server) CheckKey(c context.Context, in *pb.VertifyRequest) (*pb.VertifyResponse, error) {

	vertifyKey := in.GetKey()

	vertifyKeyToken, err := jwt.Parse(vertifyKey, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(Init.Config.Server.Vertify.SigningKey), nil
	})

	if err != nil {

		return &pb.VertifyResponse{Checked: false}, err

	}

	if vertifyKeyClaims, ok := vertifyKeyToken.Claims.(jwt.MapClaims); ok && vertifyKeyToken.Valid {

		visitCounts, ok := vertifyKeyClaims["visitCounts"].(float64)

		if !ok {
			return &pb.VertifyResponse{Checked: false}, errors.New("invalid vertify key")
		}

		if int(visitCounts) >= Init.Config.Server.Vertify.VisitCounts {
			return &pb.VertifyResponse{Checked: false}, errors.New("visit counts limit")
		}

		checkCounts, ok := vertifyKeyClaims["checkCounts"].(float64)

		if !ok {
			return &pb.VertifyResponse{Checked: false}, errors.New("invalid vertify key")
		}

		if int(checkCounts) >= Init.Config.Server.Vertify.CheckCounts {
			return &pb.VertifyResponse{Checked: false}, errors.New("check counts limit")
		}

		ip, ok := vertifyKeyClaims["iss"].(string)
		if !ok {
			return &pb.VertifyResponse{Checked: false}, errors.New("invalid vertify key")
		}

		tempVertifyKey, err := Redis.Client.Get(ip).Result()
		if err != nil {
			return &pb.VertifyResponse{Checked: false}, err
		}

		if tempVertifyKey != vertifyKey {
			return &pb.VertifyResponse{Checked: false}, errors.New("invalid vertify key")
		}

		_, err = IndexService.GenerateAndSetNewVertifyKey(ip, visitCounts, checkCounts+1)

		if err != nil {
			return &pb.VertifyResponse{Checked: false}, err
		}

		return &pb.VertifyResponse{Checked: true}, nil

	}

	return &pb.VertifyResponse{Checked: false}, errors.New("invalid vertify key")

}

func (s *Server) Start() error {

	lis, err := net.Listen("tcp", s.IP+":"+strconv.Itoa(s.Port))

	defer lis.Close()

	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer()

	pb.RegisterVertifyServer(gRPCServer, s)

	reflection.Register(gRPCServer)

	return gRPCServer.Serve(lis)

}
