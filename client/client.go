package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	pb "github.com/rahulgupta0409/authentication-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "0.0.0.0:3001", "connect server address")
)

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

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	UserType  string
	jwt.StandardClaims
}

var SECRET_KEY = ""

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if token.Claims.(*SignedDetails).ExpiresAt < time.Now().Local().Unix() {
		logger.Error(err.Error(),
			zap.String("operation", "scddssd"),
			zap.Time("timestamp", time.Now().UTC()),
		)
	}

	if err != nil {
		msg = fmt.Sprintf("your token is expired")
		logger.Error(err.Error(),
			zap.String("operation", "validation"),
			zap.Time("timestamp", time.Now().UTC()),
		)
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		logger.Info(err.Error(),
			zap.String("operation", "validation"),
			zap.Time("timestamp", time.Now().UTC()),
		)
		return
	}

	// if claims.ExpiresAt < time.Now().Local().Unix() {
	// 	msg = fmt.Sprintf("token is expired")
	// 	msg = err.Error()
	// 	return
	// }
	return claims, msg
}

// func respondWithJson(w http.ResponseWriter, code int, payload interface{}) interface{} {
// 	response, _ := json.Marshal(payload)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// 	return payload
// }

func AuthenticateWithToken() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		userId := request.URL.Query().Get("UserId")
		token := request.URL.Query().Get("Token")
		var user User

		logger := zap.Must(zap.NewProduction())
		defer logger.Sync()

		// Set up a connection to the server.
		conn, err := grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			logger.Error("connection error",
				zap.String("operation", "failed to connect with gRPC server"),
				zap.Error(errors.New(err.Error())),
				zap.Time("timestamp", time.Now().UTC()),
			)
		}
		defer conn.Close()
		grpcClient := pb.NewUserServiceClient(conn)

		stream1, err := grpcClient.AuthenticateUser(context.Background())
		if err != nil {
			logger.Error("receive stream error",
				zap.String("operation", "gRPC bidirectional stream"),
				zap.Error(errors.New(err.Error())),
				zap.Time("timestamp", time.Now().UTC()),
			)
		}
		stream1.Send(&pb.UserWithTokenRequest{UserId: userId, Token: token})

		resp, err := stream1.Recv()
		if err != nil {
			log.Fatalf("resp error: %v", err)
			logger.Error("response error",
				zap.String("operation", "gRPC bidirectional stream response"),
				zap.Error(errors.New(err.Error())),
				zap.Time("timestamp", time.Now().UTC()),
			)
		}

		user.FirstName = &resp.FirstName
		user.LastName = &resp.LastName

		result, _ := json.Marshal(resp)

		claims, msg := ValidateToken(token)
		if claims.ExpiresAt < time.Now().Local().Unix() {
			logger.Info(msg,
				zap.String("operation", "authentication"),
				zap.Time("timestamp", time.Now().UTC()),
			)
		} else {
			response.Write(result)
			logger.Info("Successfully Authenticated",
				zap.String("UserId", resp.UserId),
				zap.Time("timestamp", time.Now().UTC()),
			)
		}
		//response.Write(result)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/authenticate", AuthenticateWithToken()).Methods("POST", "OPTIONS")

	port := 8182
	fmt.Printf("Server listening on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
