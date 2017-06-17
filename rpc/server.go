package rpc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

	if vertifyKey == "" {
		return nil, errors.New("invalid vertify key")
	}

	vertifyKeyToken, err := jwt.Parse(vertifyKey, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(Init.Config.Server.Vertify.SigningKey), nil
	})

	if err != nil {

		return nil, err

	}

	if vertifyKeyClaims, ok := vertifyKeyToken.Claims.(jwt.MapClaims); ok && vertifyKeyToken.Valid {

		ip, ok := vertifyKeyClaims["iss"].(string)
		if !ok {
			return nil, errors.New("invalid vertify key")
		}

		record, err := Redis.Client.Get(ip).Result()

		if err != nil {
			return nil, err
		}

		if record == "" {
			return nil, errors.New("invalid vertify key")

		}

		if matched, err := regexp.Match(`^\d+:\d+$`, []byte(record)); !matched || err != nil {
			return nil, errors.New("invalid vertify key")
		}

		recordValues := strings.Split(record, ":")

		visitCountsValue := recordValues[0]
		visitCounts, err := strconv.Atoi(visitCountsValue)
		if err != nil {
			return nil, errors.New("invalid vertify key")
		}

		if visitCounts >= Init.Config.Server.Vertify.VisitCounts {
			return &pb.VertifyResponse{Checked: false, Status: pb.Status_VISITLIMIT}, nil
		}

		checkCountsValue := recordValues[1]
		checkCounts, err := strconv.Atoi(checkCountsValue)
		if err != nil {
			return nil, errors.New("invalid vertify key")
		}

		if checkCounts >= Init.Config.Server.Vertify.CheckCounts {
			return &pb.VertifyResponse{Checked: false, Status: pb.Status_CHECKLIMIT}, nil
		}

		_, err = IndexService.GenerateAndSetNewVertifyKey(ip, visitCounts, checkCounts+1)

		if err != nil {
			return nil, err
		}

		return &pb.VertifyResponse{Checked: true}, nil

	}

	return nil, errors.New("invalid vertify key")

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
