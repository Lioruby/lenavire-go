package queries

import (
	"encoding/json"

	"gorm.io/gorm"
)

type GetContributorsRankQuery struct{}

func NewGetContributorsRankQuery() GetContributorsRankQuery {
	return GetContributorsRankQuery{}
}

type Contributor struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	Name   string `json:"name"`
}

type GetContributorsRankQueryResult struct {
	Contributors []Contributor `json:"contributors"`
}

type GetContributorsRankQueryHandler struct {
	db *gorm.DB
}

func NewGetContributorsRankQueryHandler(db *gorm.DB) *GetContributorsRankQueryHandler {
	return &GetContributorsRankQueryHandler{db: db}
}

func (h *GetContributorsRankQueryHandler) Execute(query GetContributorsRankQuery) (*GetContributorsRankQueryResult, error) {
	var result struct {
		Contributors json.RawMessage `json:"contributors"`
	}

	sqlQuery := `
		SELECT json_agg(
			json_build_object(
				'amount', total_amount,
				'name', name
			)
		) as contributors
		FROM (
			SELECT 
				email,
				SUM(amount) as total_amount,
				MAX(name) as name
			FROM payments
			GROUP BY email
			ORDER BY total_amount DESC
			LIMIT 100
		) ranked_contributors
	`

	err := h.db.Raw(sqlQuery).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var contributors []Contributor
	if err := json.Unmarshal(result.Contributors, &contributors); err != nil {
		return nil, err
	}

	return &GetContributorsRankQueryResult{Contributors: contributors}, nil
}
