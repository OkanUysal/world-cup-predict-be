package scoring

import (
	"encoding/json"
	"fmt"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
)

type MatchScoreChoice struct {
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}

type MatchScoreResult struct {
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}

type TeamChoice struct {
	Team string `json:"team"`
}

type TeamResult struct {
	Team string `json:"team"`
}

func CalculatePoints(eventType domain.EventType, choiceJSON, resultJSON json.RawMessage) (int, error) {
	switch eventType {
	case domain.EventTypeMatchScore:
		return scoreMatch(choiceJSON, resultJSON)
	case domain.EventTypeChampion:
		return scoreTeamPick(choiceJSON, resultJSON, 10)
	case domain.EventTypeRunnerUp:
		return scoreTeamPick(choiceJSON, resultJSON, 6)
	case domain.EventTypeThirdPlace:
		return scoreTeamPick(choiceJSON, resultJSON, 4)
	default:
		return 0, fmt.Errorf("unknown event type: %s", eventType)
	}
}

func scoreMatch(choiceJSON, resultJSON json.RawMessage) (int, error) {
	var choice MatchScoreChoice
	var result MatchScoreResult

	if err := json.Unmarshal(choiceJSON, &choice); err != nil {
		return 0, fmt.Errorf("invalid choice: %w", err)
	}
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		return 0, fmt.Errorf("invalid result: %w", err)
	}

	if choice.HomeScore == result.HomeScore && choice.AwayScore == result.AwayScore {
		return 3, nil
	}

	choiceOutcome := outcome(choice.HomeScore, choice.AwayScore)
	resultOutcome := outcome(result.HomeScore, result.AwayScore)
	if choiceOutcome == resultOutcome {
		return 1, nil
	}

	return 0, nil
}

func outcome(home, away int) int {
	switch {
	case home > away:
		return 1
	case home < away:
		return -1
	default:
		return 0
	}
}

func scoreTeamPick(choiceJSON, resultJSON json.RawMessage, points int) (int, error) {
	var choice TeamChoice
	var result TeamResult

	if err := json.Unmarshal(choiceJSON, &choice); err != nil {
		return 0, fmt.Errorf("invalid choice: %w", err)
	}
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		return 0, fmt.Errorf("invalid result: %w", err)
	}

	if choice.Team == result.Team {
		return points, nil
	}

	return 0, nil
}
