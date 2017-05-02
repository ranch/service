package task

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/knq/firebase"
	"github.com/synoday/golang/auth"

	"github.com/synoday/golang/grpc/interceptor"
	pb "github.com/synoday/golang/protogen/task"
)

// Serve registers task service as gRPC server in specified connection.
func (s *Service) Serve() {
	addr := fmt.Sprintf("localhost:%d", s.config.GetInt("server.port"))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(
		interceptor.AuthUnary(),
	))
	server := grpc.NewServer(opts...)

	pb.RegisterTaskServiceServer(server, s)

	log.Printf("Task service started on: %s\n", addr)
	server.Serve(lis)
}

// TodayTasks retrieves all tasks added today.
func (s *Service) TodayTasks(ctx context.Context, req *pb.TodayTasksRequest) (*pb.TodayTasksResponse, error) {
	var err error
	res := &pb.TodayTasksResponse{
		Status: pb.ResponseStatus_INTERNAL_ERROR,
	}

	userID := ctx.Value(auth.UserIDKey).(string)
	today := time.Now().Format("20060102")
	tasks := make(map[string]Task)

	err = s.dataRef.Ref("/task/"+userID).Get(&tasks,
		firebase.OrderBy("date"),
		firebase.StartAt(today),
		firebase.EndAt(today),
	)

	if err != nil {
		return res, err
	}

	for _, task := range tasks {
		res.Tasks = append(res.Tasks, &pb.Task{
			TaskName: task.TaskName,
			Url:      task.URL,
			Tags:     task.Tags,
			Notes:    task.Notes,
			NotesMd:  task.NotesMD,
		})
	}
	res.Status = pb.ResponseStatus_SUCCESS
	return res, nil
}

// AddTask add new task to datebase.
func (s *Service) AddTask(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	var err error
	res := &pb.AddResponse{
		Status: pb.ResponseStatus_INTERNAL_ERROR,
	}

	userID := ctx.Value(auth.UserIDKey).(string)
	today := time.Now().Format("20060102")
	task := &Task{
		Date:     today,
		TaskName: req.Task.TaskName,
		URL:      req.Task.Url,
		Tags:     req.Task.Tags,
		Notes:    req.Task.Notes,
	}

	id, err := s.dataRef.Ref("/task/" + userID).Push(task)
	if err != nil {
		return res, err
	}

	return &pb.AddResponse{
		Id:     id,
		Status: pb.ResponseStatus_SUCCESS,
	}, nil
}

// RemoveTask get single task that matches with provided criteria.
func (s *Service) RemoveTask(ctx context.Context, req *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	var err error
	res := &pb.RemoveResponse{
		Status: pb.ResponseStatus_INTERNAL_ERROR,
	}

	userID := ctx.Value(auth.UserIDKey).(string)
	path := fmt.Sprintf("/task/%s/%s", userID, req.Id)

	err = s.dataRef.Ref(path).Remove()
	if err != nil {
		return res, err
	}

	return &pb.RemoveResponse{
		Status: pb.ResponseStatus_SUCCESS,
	}, nil
}
