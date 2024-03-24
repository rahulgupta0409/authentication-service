package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/rahulgupta0409/authentication-service/configs"
	pb "github.com/rahulgupta0409/authentication-service/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 3001, "listening port")

type bidirectionalStream struct {
	pb.UnimplementedUserServiceServer
}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserId       string             `json:"userid"`
	FirstName    *string            `json:"firstname"`
	LastName     *string            `json:"lastname"`
	Password     *string            `json:"password"`
	Email        *string            `json:"email"`
	Phone        *string            `json:"phone"`
	IsActive     bool               `json:"isactive"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshtoken"`
	UserType     *string            `json:"usertype"`
	CreatedBy    *string            `json:"createdby"`
	CreatedDate  time.Time          `json:"createddate"`
	ModifiedBy   *string            `json:"modifiedby"`
	ModifiedDate time.Time          `json:"modifieddate"`
}

func (s *bidirectionalStream) AuthenticateUser(stream pb.UserService_AuthenticateUserServer) error {
	var foundUser User
	var ctx, cancel = context.WithTimeout(context.Background(), 1000000*time.Hour)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = userCollection.FindOne(ctx, bson.M{"userid": req.UserId}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			return err
		}

		if req.Token == *foundUser.Token {
			u := pb.UserResponse{UserId: foundUser.UserId, FirstName: *foundUser.FirstName,
				LastName: *foundUser.LastName, Password: *foundUser.Password,
				Email: *foundUser.Email}

			err = stream.Send(&u)
			if err != nil {
				return err
			}
			log.Printf("[RECEVIED REQUEST] : %v\n", req)

		}

	}
}

func main() {
	var s *bidirectionalStream
	flag.Parse()
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listen error: %v", err)
	} else {
		log.Printf("server listen: ", addr)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, s)
	grpcServer.Serve(listener)

}
