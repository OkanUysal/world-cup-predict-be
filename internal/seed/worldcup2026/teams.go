package worldcup2026

import "time"

// Teams lists all 48 nations at FIFA World Cup 2026 (Turkish names).
var Teams = []string{
	"ABD",
	"Almanya",
	"Arjantin",
	"Avustralya",
	"Avusturya",
	"Belçika",
	"Bosna-Hersek",
	"Brezilya",
	"Cezayir",
	"Curaçao",
	"Çekya",
	"Ekvador",
	"Fas",
	"Fildişi Sahili",
	"Fransa",
	"Gana",
	"Güney Afrika",
	"Güney Kore",
	"Haiti",
	"Hırvatistan",
	"Hollanda",
	"İngiltere",
	"İran",
	"Irak",
	"İskoçya",
	"İspanya",
	"İsveç",
	"İsviçre",
	"Japonya",
	"Kanada",
	"Katar",
	"Kolombiya",
	"Kongo DC",
	"Meksika",
	"Mısır",
	"Norveç",
	"Panama",
	"Paraguay",
	"Portekiz",
	"Senegal",
	"Suudi Arabistan",
	"Tunus",
	"Türkiye",
	"Uruguay",
	"Ürdün",
	"Yeşil Burun Adaları",
	"Yeni Zelanda",
	"Özbekistan",
}

// TournamentOpens is the opening match kickoff in UTC (GMT).
var TournamentOpens = time.Date(2026, 6, 11, 19, 0, 0, 0, time.UTC)

// PredictionDeadline is one day before the tournament opens (UTC/GMT).
var PredictionDeadline = TournamentOpens.Add(-24 * time.Hour)

const SeedTitlePrefix = "WC2026:"
