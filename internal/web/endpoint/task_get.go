package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/owlint/maestro/internal/domain"
	taskerrors "github.com/owlint/maestro/internal/errors"
	"github.com/owlint/maestro/internal/infrastructure/persistence/view"
)

// TaskStateRequest is the request for the state of a task
type TaskStateRequest struct {
	TaskID string `json:"task_id"`
}

type TaskDTO struct {
	TaskID     string `json:"task_id"`
	Owner      string `json:"owner"`
	TaskQueue  string `json:"task_queue"`
	Payload    string `json:"payload"`
	State      string `json:"state"`
	Timeout    int32  `json:"timeout"`
	Retries    int32  `json:"retries"`
	MaxRetries int32  `json:"max_retries"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	Result     string `json:"result,omitempty"`
	NotBefore  int64  `json:"not_before"`
}

// TaskStateResponse is the response of a task state
type TaskStateResponse struct {
	Task  *TaskDTO `json:"task,omitempty"`
	Error string   `json:"error,omitempty"`
}

// TaskStateEndpoint creates a endpoint for task state
func TaskStateEndpoint(svc view.TaskView) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := unmarshalTaskStateRequest(request)
		if err != nil {
			return TaskStateResponse{nil, err.Error()}, taskerrors.ValidationError{Origin: err}
		}

		task, err := svc.ByID(ctx, req.TaskID)
		if err != nil {
			return TaskStateResponse{nil, err.Error()}, err
		}

		taskDTO := fromTask(task)
		return TaskStateResponse{&taskDTO, ""}, nil
	}
}

func fromTask(task *domain.Task) TaskDTO {
	var result string
	if task.State() == domain.TaskStateCompleted {
		result, _ = task.Result()
	}
	return TaskDTO{
		TaskID:     task.TaskID,
		Owner:      task.Owner(),
		TaskQueue:  task.Queue(),
		Payload:    task.Payload(),
		State:      task.State().String(),
		Timeout:    task.GetTimeout(),
		Retries:    task.Retries(),
		MaxRetries: task.MaxRetries(),
		CreatedAt:  task.CreatedAt(),
		UpdatedAt:  task.UpdatedAt(),
		Result:     result,
		NotBefore:  task.NotBefore(),
	}
}

func unmarshalTaskStateRequest(request interface{}) (*TaskStateRequest, error) {
	req := request.(TaskStateRequest)
	if req.TaskID == "" {
		return nil, errors.New("task_id is a required parameter")
	}
	return &req, nil
}
