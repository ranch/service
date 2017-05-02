package task

import (
	"log"
	"os"
	"testing"

	"golang.org/x/net/context"

	"github.com/knq/firebase"
	"github.com/synoday/synoday/golang/auth"
	pb "github.com/synoday/synoday/golang/protogen/task"
)

var testService *Service
var defaultContext context.Context
var defaultCancel context.CancelFunc
var userID string

func mockData() {
	var req pb.AddRequest

	testData := []pb.Task{
		{
			TaskName: "Foo",
			Tags:     "Bar",
			Notes:    "# Baz",
			Url:      "brank.as",
		},
		{
			TaskName: "Kudu",
			Tags:     "App",
			Notes:    "## Test",
			Url:      "google.com",
		},
	}

	for _, test := range testData {
		req.Task = &test
		_, err := testService.AddTask(defaultContext, &req)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func clearData() {
	var err error

	err = testService.dataRef.Ref("/task").Remove()
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	testService = NewService(
		ConfigFile("./testdata/env/config"),
	)

	//TODO: Replace dummy user ID.
	userID = "foo"
	defaultContext, defaultCancel = context.WithCancel(context.WithValue(context.Background(), auth.UserIDKey, userID))

	clearData()
	mockData()

	code := m.Run()

	os.Exit(code)
}

func TestTodayItems(t *testing.T) {
	var err error

	test := pb.TodayTasksRequest{
		TaskName: "Foo",
		Tags:     "Bar",
	}
	res, err := testService.TodayTasks(defaultContext, &test)
	if err != nil {
		t.Fatal(err)
	}

	if res.Status != pb.ResponseStatus_SUCCESS {
		t.Fatalf("expected response status is: '%v', got: '%v'",
			pb.ResponseStatus_SUCCESS,
			res.Status)
	}
}

func TestAddItem(t *testing.T) {
	var err error
	var req pb.AddRequest

	req.Task = &pb.Task{
		TaskName: "Foo",
		Tags:     "Bar",
		Notes:    "# Baz",
		Url:      "reddit.com",
	}
	res, err := testService.AddTask(defaultContext, &req)
	if err != nil {
		t.Fatal(err)
	}

	if res.Status != pb.ResponseStatus_SUCCESS {
		t.Fatalf("expected response status is: '%v', got: '%v'",
			pb.ResponseStatus_SUCCESS,
			res.Status)
	}

	if res.Id == "" {
		t.Error("Expected id to not empty")
	}
}

func TestRemoveItem(t *testing.T) {
	var err error

	keys := make(map[string]interface{})
	err = testService.dataRef.Ref("/task/"+userID).Get(&keys, firebase.Shallow)
	if err != nil {
		log.Fatal(err)
	}

	if len(keys) < 1 {
		log.Fatalf("expected at least one task to be present")
	}

	for key := range keys {
		res, err := testService.RemoveTask(defaultContext, &pb.RemoveRequest{Id: key})
		if err != nil {
			t.Fatal(err)
		}

		if res.Status != pb.ResponseStatus_SUCCESS {
			t.Fatalf("expected response status is: '%v', got: '%v'",
				pb.ResponseStatus_SUCCESS,
				res.Status)
		}
	}
}
