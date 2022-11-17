package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotExists = errors.New("Row does not exists")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) AggregatedCategoryScoresForPeriod(start_date, end_date int64) ([]CategoryScore, error) {
	q := `
	select
	rating_categories.name as Category,
	group_concat(
	(case
		when (julianday(datetime(?, 'unixepoch')) - julianday(datetime(?, 'unixepoch'))) > 30 then (strftime('%W', strftime('%Y-%m-%d', ratings.created_at)))
		else (strftime('%Y-%m-%d', ratings.created_at))
	end) || ',' || (case
			when (ratings.rating * rating_categories.weight / 5.0 * 100) > 100 then 100
			else ratings.rating * rating_categories.weight / 5.0 * 100
		end),';'
	) as ScoreDate,
	count(ratings.id) as Rating,
	round(avg(ratings.rating) * rating_categories.weight / 5 * 100,
	2) as score
	from
		ratings
	join rating_categories on
		rating_categories.id = ratings.rating_category_id
	where
		ratings.created_at between datetime(?, 'unixepoch') and datetime(?, 'unixepoch')
	group by
		rating_categories.name;
		`
	rows, err := r.db.Query(q, end_date, start_date, start_date, end_date, end_date, start_date)

	if err != nil {
		return nil, fmt.Errorf("Couldn't query data, %w", err)
	}
	defer rows.Close()

	var all []CategoryScore
	for rows.Next() {
		var cs CategoryScore
		if err := rows.Scan(&cs.Category, &cs.DateScores, &cs.RatingCount, &cs.TotalScore); err != nil {
			return nil, fmt.Errorf("Failed to scan data into struct, %w", err)
		}
		all = append(all, cs)
	}
	return all, nil
}

func (r *SQLiteRepository) TicketScoresForPeriod(start_date, end_date int64) ([]TicketScore, error) {
	q := `
	select
	tickets.id,
	group_concat(
		rating_categories.name || ',' || (case
			when (ratings.rating * rating_categories.weight / 5.0 * 100) > 100 then 100
			else ratings.rating * rating_categories.weight / 5.0 * 100
		end),
	';'
	) as category_scores
	from
		ratings
	join rating_categories on
		rating_categories.id = ratings.rating_category_id
	join tickets on
		tickets.id = ratings.ticket_id
	where
		ratings.created_at BETWEEN datetime(?, 'unixepoch') AND datetime(?, 'unixepoch')
	group by
		tickets.id;
		`
	rows, err := r.db.Query(q, start_date, end_date)

	if err != nil {
		return nil, fmt.Errorf("Failed to query rows, %w", err)
	}
	defer rows.Close()

	var all []TicketScore
	for rows.Next() {
		var ts TicketScore
		if err := rows.Scan(&ts.ID, &ts.CategoryScores); err != nil {
			return nil, fmt.Errorf("Failed to scan data into struct, %w", err)
		}
		all = append(all, ts)
	}
	return all, nil
}

func (r *SQLiteRepository) OverallScoresForPeriod(start_date, end_date int64) (Score, error) {
	q := `
	select
	round(avg(case
		when (ratings.rating * rating_categories.weight / 5.0 * 100) > 100 then 100
		else ratings.rating * rating_categories.weight / 5.0 * 100
	end),2) as overal_score
	from
		ratings
	join rating_categories on
		rating_categories.id = ratings.rating_category_id
	where
		ratings.created_at between datetime(?, 'unixepoch') and datetime(?, 'unixepoch');
		`
	row := r.db.QueryRow(q, start_date, end_date)
	var score Score
	if err := row.Scan(&score.Score); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return score, ErrNotExists
		}
		return score, err
	}
	return score, nil
}
