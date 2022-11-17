package presenter

import (
	"strconv"
	"strings"

	pb "github.com/valdrinkuchi/score_server/proto_files"
	"github.com/valdrinkuchi/score_server/repository"
)

func TicketPresenter(data []repository.TicketScore) *pb.TicketScoresResponse {
	var result []*pb.TicketScore
	for _, scores := range data {
		result = append(result, &pb.TicketScore{Id: scores.ID, CategoryScores: categoryScoreMapper(scores.CategoryScores)})
	}
	res := &pb.TicketScoresResponse{
		TicketScores: result,
	}
	return res
}

func categoryScoreMapper(score string) []*pb.CategoryScore {
	splitedScores := strings.Split(score, ";")
	var cs []*pb.CategoryScore
	for _, score := range splitedScores {
		splittedScore := strings.Split(score, ",")
		cs = append(cs, &pb.CategoryScore{CategoryName: splittedScore[0], Score: stringToFloat(splittedScore[1])})
	}
	return cs
}

func stringToFloat(value string) float64 {
	v, _ := strconv.ParseFloat(value, 64)
	return v
}
