package scoring

import (
	"encoding/json"
	"testing"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
)

func TestScoreMatchExact(t *testing.T) {
	choice := json.RawMessage(`{"home_score":2,"away_score":1}`)
	result := json.RawMessage(`{"home_score":2,"away_score":1}`)

	points, err := CalculatePoints(domain.EventTypeMatchScore, choice, result)
	if err != nil {
		t.Fatal(err)
	}
	if points != 3 {
		t.Fatalf("expected 3 points, got %d", points)
	}
}

func TestScoreMatchOutcome(t *testing.T) {
	choice := json.RawMessage(`{"home_score":3,"away_score":0}`)
	result := json.RawMessage(`{"home_score":2,"away_score":1}`)

	points, err := CalculatePoints(domain.EventTypeMatchScore, choice, result)
	if err != nil {
		t.Fatal(err)
	}
	if points != 1 {
		t.Fatalf("expected 1 point, got %d", points)
	}
}

func TestScoreMatchWrong(t *testing.T) {
	choice := json.RawMessage(`{"home_score":0,"away_score":1}`)
	result := json.RawMessage(`{"home_score":2,"away_score":1}`)

	points, err := CalculatePoints(domain.EventTypeMatchScore, choice, result)
	if err != nil {
		t.Fatal(err)
	}
	if points != 0 {
		t.Fatalf("expected 0 points, got %d", points)
	}
}

func TestScoreChampion(t *testing.T) {
	choice := json.RawMessage(`{"team":"Brazil"}`)
	result := json.RawMessage(`{"team":"Brazil"}`)

	points, err := CalculatePoints(domain.EventTypeChampion, choice, result)
	if err != nil {
		t.Fatal(err)
	}
	if points != 10 {
		t.Fatalf("expected 10 points, got %d", points)
	}
}

func TestScoreRunnerUp(t *testing.T) {
	choice := json.RawMessage(`{"team":"Argentina"}`)
	result := json.RawMessage(`{"team":"Argentina"}`)

	points, err := CalculatePoints(domain.EventTypeRunnerUp, choice, result)
	if err != nil {
		t.Fatal(err)
	}
	if points != 5 {
		t.Fatalf("expected 5 points, got %d", points)
	}
}
