package worldcup2026

import "time"

// Teams lists all 48 nations at FIFA World Cup 2026.
var Teams = []string{
	"Algeria",
	"Argentina",
	"Australia",
	"Austria",
	"Belgium",
	"Bosnia and Herzegovina",
	"Brazil",
	"Cabo Verde",
	"Canada",
	"Colombia",
	"Congo DR",
	"Croatia",
	"Curaçao",
	"Czechia",
	"Côte d'Ivoire",
	"Ecuador",
	"Egypt",
	"England",
	"France",
	"Germany",
	"Ghana",
	"Haiti",
	"IR Iran",
	"Iraq",
	"Japan",
	"Jordan",
	"Korea Republic",
	"Mexico",
	"Morocco",
	"Netherlands",
	"New Zealand",
	"Norway",
	"Panama",
	"Paraguay",
	"Portugal",
	"Qatar",
	"Saudi Arabia",
	"Scotland",
	"Senegal",
	"South Africa",
	"Spain",
	"Sweden",
	"Switzerland",
	"Tunisia",
	"Türkiye",
	"USA",
	"Uruguay",
	"Uzbekistan",
}

// TournamentOpens is the opening match kickoff in UTC (GMT).
var TournamentOpens = time.Date(2026, 6, 11, 19, 0, 0, 0, time.UTC)

// PredictionDeadline is one day before the tournament opens (UTC/GMT).
var PredictionDeadline = TournamentOpens.Add(-24 * time.Hour)

const SeedTitlePrefix = "WC2026:"
