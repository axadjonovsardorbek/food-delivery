package service

import (
	"context"
	cp "courier/genproto/courier"
	st "courier/storage/postgres"
)

type TaskService struct {
	storage st.Storage
	cp.UnimplementedTaskServiceServer
}

func NewTaskService(storage *st.Storage) *TaskService {
	return &TaskService{storage: *storage}
}

func (s *TaskService) Create(ctx context.Context, req *cp.TaskCreateReq) (*cp.Void, error) {
	return s.storage.TaskS.Create(req)
}
func (s *TaskService) GetById(ctx context.Context, req *cp.ById) (*cp.TaskGetByIdRes, error) {
	return s.storage.TaskS.GetById(req)
}
func (s *TaskService) GetAll(ctx context.Context, req *cp.TaskGetAllReq) (*cp.TaskGetAllRes, error) {
	return s.storage.TaskS.GetAll(req)
}
func (s *TaskService) Update(ctx context.Context, req *cp.TaskUpdateReq) (*cp.Void, error) {
	return s.storage.TaskS.Update(req)
}
func (s *TaskService) Delete(ctx context.Context, req *cp.ById) (*cp.Void, error) {
	return s.storage.TaskS.Delete(req)
}
