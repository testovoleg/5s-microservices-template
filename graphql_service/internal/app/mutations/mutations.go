package mutations

import (
	"time"

	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Mutations struct {
	CreateBug CreateBugHandler
}

func NewMutations(
	createBug CreateBugHandler,
) *Mutations {
	return &Mutations{
		CreateBug: createBug,
	}
}

type CreateBugCommand struct {
	AccessToken        string
	Name               string
	Description        string
	CreateForReleaseID string
	Files              []string
}

func NewCreateBugCommand(access_token string, in *model.NewBug) *CreateBugCommand {
	return &CreateBugCommand{
		AccessToken:        access_token,
		Name:               in.Name,
		Description:        DerefString(in.Description),
		CreateForReleaseID: in.CreateForReleaseID,
		Files:              in.Files,
	}
}

func StrToTimestamppb(str *string) (*timestamppb.Timestamp, error) {
	if str == nil || *str == "" {
		return nil, nil
	}
	t, err := time.Parse("02.01.2006", *str)
	if err != nil {
		return nil, err
	}

	return timestamppb.New(t), nil
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
