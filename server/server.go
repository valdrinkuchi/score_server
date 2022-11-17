package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/valdrinkuchi/score_server/presenter"
	pb "github.com/valdrinkuchi/score_server/proto_files"
	"github.com/valdrinkuchi/score_server/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ADDRESS = "localhost:50051"
)

type server struct {
	pb.UnimplementedScoresServer
}

func (*server) GetAggregatedCategoryScoresForPeriod(ctx context.Context, req *pb.Interval) (*pb.AggregatedCategoryScoresResponse, error) {
	fmt.Println("GetAggregatedCategoryScoresForPeriod Called")
	start_date := req.GetStartDate()
	end_date := req.GetEndDate()
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "cannot open sqlite database")
	}

	repository := repository.NewSQLiteRepository(db)

	data, err := repository.AggregatedCategoryScoresForPeriod(start_date, end_date)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return presenter.ScoresPresenter(data), nil
}

func (*server) GetTicketScoresForPeriod(ctx context.Context, req *pb.Interval) (*pb.TicketScoresResponse, error) {
	fmt.Println("GetTicketScoreForPeriod Called")
	start_date := req.GetStartDate()
	end_date := req.GetEndDate()

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "cannot open sqlite database")
	}

	repository := repository.NewSQLiteRepository(db)

	data, err := repository.TicketScoresForPeriod(start_date, end_date)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return presenter.TicketPresenter(data), nil
}

func (*server) GetOverallScoreForPeriod(ctx context.Context, req *pb.Interval) (*pb.OverallScoreResponse, error) {
	fmt.Println("GetOverallScoreForPeriod Called")
	start_date := req.GetStartDate()
	end_date := req.GetEndDate()
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "cannot open sqlite database")
	}

	repository := repository.NewSQLiteRepository(db)

	data, err := repository.OveralScoresForPeriod(start_date, end_date)
	if err != nil {
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return presenter.OverallScorePresenter(data), nil
}

func main() {
	fmt.Println("Boom")
	lis, err := net.Listen("tcp", ADDRESS)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterScoresServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
