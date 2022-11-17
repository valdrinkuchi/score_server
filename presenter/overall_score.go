package presenter

import (
	pb "github.com/valdrinkuchi/score_server/proto_files"
	"github.com/valdrinkuchi/score_server/repository"
)

func OverallScorePresenter(data repository.Score) *pb.OverallScoreResponse {
	return &pb.OverallScoreResponse{
		Score: data.Score,
	}
}
