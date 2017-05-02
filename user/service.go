package user

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"github.com/synoday/golang/auth/token"
	pb "github.com/synoday/golang/protogen/user"
)

var (
	// ErrInvalidCredential is the error returned when the credentials
	// is not match.
	ErrInvalidCredential = errors.New("Invalid username or password")

	// ErrRegisteredEmail is the error returned when user register
	// with registered email.
	ErrRegisteredEmail = errors.New("Email already registered")

	// ErrRegisteredUsername is the error returned when user register
	// with registered email.
	ErrRegisteredUsername = errors.New("Username already registered")
)

// Serve registers user service as gRPC server in specified connection.
func (s *Service) Serve() {
	addr := fmt.Sprintf("localhost:%d", s.config.GetInt("server.port"))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, s)

	log.Printf("User service started on: %s\n", addr)
	server.Serve(lis)
}

// Register set new user record into user database ref.
func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var err error
	res := &pb.RegisterResponse{
		Status: pb.ResponseStatus_INTERNAL_ERROR,
	}

	encEmail := base64.StdEncoding.EncodeToString([]byte(req.Credential.Email))
	encUsername := base64.StdEncoding.EncodeToString([]byte(req.Credential.Username))

	// check if email / username already registered.
	var check struct {
		UserID string `json:"user_id"`
	}
	err = s.authRef.Ref("/credential/email/" + encEmail).Get(&check)
	if err != nil {
		return res, err
	}
	if check.UserID != "" {
		return res, ErrRegisteredEmail
	}

	err = s.authRef.Ref("/credential/username/" + encUsername).Get(&check)
	if err != nil {
		return res, err
	}
	if check.UserID != "" {
		return res, ErrRegisteredUsername
	}

	// Create user auth record
	userID, err := s.authRef.Ref("/user").Push(map[string]interface{}{
		"credentials": map[string]interface{}{
			"email": map[string]interface{}{
				encEmail: true,
			},
			"username": map[string]interface{}{
				encUsername: true,
			},
		},
	})
	if err != nil {
		return res, err
	}

	// Hash password
	passBuf, err := bcrypt.GenerateFromPassword([]byte(req.Credential.Password), bcrypt.DefaultCost)
	if err != nil {
		return res, err
	}

	// Set email credential
	err = s.authRef.Ref("/credential/email").Update(map[string]interface{}{
		encEmail: map[string]interface{}{
			"secret":  string(passBuf),
			"user_id": userID,
		},
	})
	if err != nil {
		return res, err
	}

	// Set username credential
	err = s.authRef.Ref("/credential/username").Update(map[string]interface{}{
		encUsername: map[string]interface{}{
			"secret":  string(passBuf),
			"user_id": userID,
		},
	})
	if err != nil {
		return res, err
	}

	// Create user data record
	user := &User{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
	}
	err = s.dataRef.Ref("/user/" + userID).Set(user)
	if err != nil {
		return res, err
	}

	return &pb.RegisterResponse{
		Status: pb.ResponseStatus_SUCCESS,
	}, nil
}

// Login validates user claims and generate jwt token.
func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var err error
	res := &pb.LoginResponse{
		Status: pb.ResponseStatus_INTERNAL_ERROR,
	}

	// TODO: Handle login with username

	// find email / username
	encEmail := base64.StdEncoding.EncodeToString([]byte(req.Credential.Email))
	var creds struct {
		UserID string `json:"user_id"`
		Secret string `json:"secret"`
	}
	err = s.authRef.Ref("/credential/email/" + encEmail).Get(&creds)
	if err != nil {
		return res, err
	}

	if creds.UserID == "" {
		res.Status = pb.ResponseStatus_CREDENTIAL_INVALID
		return res, ErrInvalidCredential
	}

	// validate credential
	err = bcrypt.CompareHashAndPassword([]byte(creds.Secret), []byte(req.Credential.Password))
	if err != nil {
		res.Status = pb.ResponseStatus_CREDENTIAL_INVALID
		return res, ErrInvalidCredential
	}

	var user User
	err = s.dataRef.Ref("/user/" + creds.UserID).Get(&user)
	if err != nil {
		return res, err
	}

	// generate token
	exp := time.Now().Add(token.DefaultExp)
	claims := &token.Claims{
		UID:        creds.UserID,
		IssuedAt:   json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		Expiration: json.Number(strconv.FormatInt(exp.Unix(), 10)),
	}
	token, err := token.New(claims,
		s.config.GetKey("jwt.privatekey"),
		s.config.GetKey("jwt.publickey"),
	)
	if err != nil {
		return res, err
	}

	return &pb.LoginResponse{
		Status: pb.ResponseStatus_SUCCESS,
		Token:  token,
	}, nil
}
