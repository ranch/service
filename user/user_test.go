package user

import (
	"context"
	"log"
	"os"
	"testing"

	pcreds "github.com/synoday/synoday/golang/protogen/type/creds"
	pb "github.com/synoday/synoday/golang/protogen/user"
)

var testService *Service
var defaultContext, defaultCancel = context.WithCancel(context.Background())

func mockData() {
	req := &pb.RegisterRequest{
		Credential: &pcreds.Credential{
			Email:    "foo@bar.com",
			Username: "foo bar baz",
			Password: "foobar",
		},
		User: &pb.User{
			FirstName: "Foo",
			LastName:  "Bar",
		},
	}

	_, err := testService.Register(defaultContext, req)
	if err != nil {
		log.Fatalf("got error: %v", err)
	}
}

func clearData() {
	var err error

	// remove user data entries
	err = testService.dataRef.Ref("/user").Remove()
	if err != nil {
		log.Fatal(err)
	}
	// remove auth user entries
	err = testService.authRef.Ref("/user").Remove()
	if err != nil {
		log.Fatal(err)
	}
	// remove auth credential entries
	err = testService.authRef.Ref("/credential").Remove()
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	testService = NewService(
		ConfigFile("./testdata/env/config"),
	)

	clearData()
	mockData()
	code := m.Run()

	os.Exit(code)
}

func TestRegister(t *testing.T) {
	req := &pb.RegisterRequest{
		Credential: &pcreds.Credential{
			Email:    "john.doe@foobar.com",
			Username: "johndoe",
			Password: "foobar",
		},
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
		},
	}

	res, err := testService.Register(defaultContext, req)
	if err != nil {
		t.Fatal(err)
	}

	if res.Status != pb.ResponseStatus_SUCCESS {
		t.Fatalf("expected response status is: '%v', got: '%v'",
			pb.ResponseStatus_SUCCESS,
			res.Status)
	}
}

func TestLoginWithEmail(t *testing.T) {
	req := &pb.LoginRequest{
		Credential: &pcreds.Credential{
			Email:    "foo@bar.com",
			Password: "foobar",
		},
	}
	res, err := testService.Login(defaultContext, req)
	if err != nil {
		t.Fatal(err)
	}

	if res.Status != pb.ResponseStatus_SUCCESS {
		t.Fatalf("expected response status is: '%v', got: '%v'",
			pb.ResponseStatus_SUCCESS,
			res.Status)
	}
}
