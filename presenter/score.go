package presenter

import (
	"strings"

	pb "github.com/valdrinkuchi/score_server/proto_files"
	"github.com/valdrinkuchi/score_server/repository"
)

func ScoresPresenter(data []repository.CategoryScore) *pb.AggregatedCategoryScoresResponse {
	var result []*pb.AggregatedCategoryScore
	for _, score := range data {
		result = append(result, &pb.AggregatedCategoryScore{CategoryName: score.Category, RatingCount: score.RatingCount, Dates: dateScoreMapper(score.DateScores), Score: score.TotalScore})
	}
	res := &pb.AggregatedCategoryScoresResponse{
		Data: result,
	}
	return res
}

func dateScoreMapper(data string) []*pb.DateScores {
	splitedData := strings.Split(data, ";")
	var cs []*pb.DateScores
	for _, score := range splitedData {
		splitedDataByComma := strings.Split(score, ",")
		cs = append(cs, &pb.DateScores{Date: splitedDataByComma[0], Score: stringToFloat(splitedDataByComma[1])})
	}
	return cs
}
