package worldcup2026

import "time"

type GroupMatch struct {
	MatchNo  int
	Group    string
	HomeTeam string
	AwayTeam string
	Kickoff  time.Time
	Venue    string
}

// GroupMatches contains all 72 group-stage fixtures.
// Kickoff times are UTC (GMT). Deadlines for predictions are kickoff minus 1 hour.
var GroupMatches = []GroupMatch{
	{1, "A", "Meksika", "Güney Afrika", utc(2026, 6, 11, 19, 0), "Mexico City"},
	{2, "A", "Güney Kore", "Çekya", utc(2026, 6, 12, 2, 0), "Guadalajara"},
	{3, "B", "Kanada", "Bosna-Hersek", utc(2026, 6, 12, 19, 0), "Toronto"},
	{4, "D", "ABD", "Paraguay", utc(2026, 6, 13, 1, 0), "Los Angeles"},
	{8, "B", "Katar", "İsviçre", utc(2026, 6, 13, 19, 0), "San Francisco Bay Area"},
	{7, "C", "Brezilya", "Fas", utc(2026, 6, 13, 22, 0), "New York/New Jersey"},
	{5, "C", "Haiti", "İskoçya", utc(2026, 6, 14, 1, 0), "Boston"},
	{6, "D", "Avustralya", "Türkiye", utc(2026, 6, 14, 4, 0), "Vancouver"},
	{10, "E", "Almanya", "Curaçao", utc(2026, 6, 14, 17, 0), "Houston"},
	{11, "F", "Hollanda", "Japonya", utc(2026, 6, 14, 20, 0), "Dallas"},
	{9, "E", "Fildişi Sahili", "Ekvador", utc(2026, 6, 14, 23, 0), "Philadelphia"},
	{12, "F", "İsveç", "Tunus", utc(2026, 6, 15, 2, 0), "Monterrey"},
	{14, "H", "İspanya", "Yeşil Burun Adaları", utc(2026, 6, 15, 16, 0), "Atlanta"},
	{16, "G", "Belçika", "Mısır", utc(2026, 6, 15, 19, 0), "Seattle"},
	{13, "H", "Suudi Arabistan", "Uruguay", utc(2026, 6, 15, 22, 0), "Miami"},
	{15, "G", "İran", "Yeni Zelanda", utc(2026, 6, 16, 1, 0), "Los Angeles"},
	{17, "I", "Fransa", "Senegal", utc(2026, 6, 16, 19, 0), "New York/New Jersey"},
	{18, "I", "Irak", "Norveç", utc(2026, 6, 16, 22, 0), "Boston"},
	{19, "J", "Arjantin", "Cezayir", utc(2026, 6, 17, 1, 0), "Kansas City"},
	{20, "J", "Avusturya", "Ürdün", utc(2026, 6, 17, 4, 0), "San Francisco Bay Area"},
	{23, "K", "Portekiz", "Kongo DC", utc(2026, 6, 17, 17, 0), "Houston"},
	{22, "L", "İngiltere", "Hırvatistan", utc(2026, 6, 17, 20, 0), "Dallas"},
	{21, "L", "Gana", "Panama", utc(2026, 6, 17, 23, 0), "Toronto"},
	{24, "K", "Özbekistan", "Kolombiya", utc(2026, 6, 18, 2, 0), "Mexico City"},
	{25, "A", "Çekya", "Güney Afrika", utc(2026, 6, 18, 16, 0), "Atlanta"},
	{26, "B", "İsviçre", "Bosna-Hersek", utc(2026, 6, 18, 19, 0), "Los Angeles"},
	{27, "B", "Kanada", "Katar", utc(2026, 6, 18, 22, 0), "Vancouver"},
	{28, "A", "Meksika", "Güney Kore", utc(2026, 6, 19, 1, 0), "Guadalajara"},
	{32, "D", "ABD", "Avustralya", utc(2026, 6, 19, 19, 0), "Seattle"},
	{30, "C", "İskoçya", "Fas", utc(2026, 6, 19, 22, 0), "Boston"},
	{29, "C", "Brezilya", "Haiti", utc(2026, 6, 20, 0, 30), "Philadelphia"},
	{31, "D", "Türkiye", "Paraguay", utc(2026, 6, 20, 3, 0), "San Francisco Bay Area"},
	{35, "F", "Hollanda", "İsveç", utc(2026, 6, 20, 17, 0), "Houston"},
	{33, "E", "Almanya", "Fildişi Sahili", utc(2026, 6, 20, 20, 0), "Toronto"},
	{34, "E", "Ekvador", "Curaçao", utc(2026, 6, 21, 0, 0), "Kansas City"},
	{36, "F", "Tunus", "Japonya", utc(2026, 6, 21, 4, 0), "Monterrey"},
	{38, "H", "İspanya", "Suudi Arabistan", utc(2026, 6, 21, 16, 0), "Atlanta"},
	{39, "G", "Belçika", "İran", utc(2026, 6, 21, 19, 0), "Los Angeles"},
	{37, "H", "Uruguay", "Yeşil Burun Adaları", utc(2026, 6, 21, 22, 0), "Miami"},
	{40, "G", "Yeni Zelanda", "Mısır", utc(2026, 6, 22, 1, 0), "Vancouver"},
	{43, "J", "Arjantin", "Avusturya", utc(2026, 6, 22, 17, 0), "Dallas"},
	{42, "I", "Fransa", "Irak", utc(2026, 6, 22, 21, 0), "Philadelphia"},
	{41, "I", "Norveç", "Senegal", utc(2026, 6, 23, 0, 0), "New York/New Jersey"},
	{44, "J", "Ürdün", "Cezayir", utc(2026, 6, 23, 3, 0), "San Francisco Bay Area"},
	{47, "K", "Portekiz", "Özbekistan", utc(2026, 6, 23, 17, 0), "Houston"},
	{45, "L", "İngiltere", "Gana", utc(2026, 6, 23, 20, 0), "Boston"},
	{46, "L", "Panama", "Hırvatistan", utc(2026, 6, 23, 23, 0), "Toronto"},
	{48, "K", "Kolombiya", "Kongo DC", utc(2026, 6, 24, 2, 0), "Guadalajara"},
	{51, "B", "İsviçre", "Kanada", utc(2026, 6, 24, 19, 0), "Vancouver"},
	{52, "B", "Bosna-Hersek", "Katar", utc(2026, 6, 24, 19, 0), "Seattle"},
	{49, "C", "İskoçya", "Brezilya", utc(2026, 6, 24, 22, 0), "Miami"},
	{50, "C", "Fas", "Haiti", utc(2026, 6, 24, 22, 0), "Atlanta"},
	{53, "A", "Çekya", "Meksika", utc(2026, 6, 25, 1, 0), "Mexico City"},
	{54, "A", "Güney Afrika", "Güney Kore", utc(2026, 6, 25, 1, 0), "Monterrey"},
	{55, "E", "Curaçao", "Fildişi Sahili", utc(2026, 6, 25, 20, 0), "Philadelphia"},
	{56, "E", "Ekvador", "Almanya", utc(2026, 6, 25, 20, 0), "New York/New Jersey"},
	{57, "F", "Japonya", "İsveç", utc(2026, 6, 25, 23, 0), "Dallas"},
	{58, "F", "Tunus", "Hollanda", utc(2026, 6, 25, 23, 0), "Kansas City"},
	{59, "D", "Türkiye", "ABD", utc(2026, 6, 26, 2, 0), "Los Angeles"},
	{60, "D", "Paraguay", "Avustralya", utc(2026, 6, 26, 2, 0), "San Francisco Bay Area"},
	{61, "I", "Norveç", "Fransa", utc(2026, 6, 26, 19, 0), "Boston"},
	{62, "I", "Senegal", "Irak", utc(2026, 6, 26, 19, 0), "Toronto"},
	{65, "H", "Yeşil Burun Adaları", "Suudi Arabistan", utc(2026, 6, 27, 0, 0), "Houston"},
	{66, "H", "Uruguay", "İspanya", utc(2026, 6, 27, 0, 0), "Guadalajara"},
	{63, "G", "Mısır", "İran", utc(2026, 6, 27, 3, 0), "Seattle"},
	{64, "G", "Yeni Zelanda", "Belçika", utc(2026, 6, 27, 3, 0), "Vancouver"},
	{67, "L", "Panama", "İngiltere", utc(2026, 6, 27, 21, 0), "New York/New Jersey"},
	{68, "L", "Hırvatistan", "Gana", utc(2026, 6, 27, 21, 0), "Philadelphia"},
	{71, "K", "Kolombiya", "Portekiz", utc(2026, 6, 27, 23, 30), "Miami"},
	{72, "K", "Kongo DC", "Özbekistan", utc(2026, 6, 27, 23, 30), "Atlanta"},
	{69, "J", "Cezayir", "Avusturya", utc(2026, 6, 28, 2, 0), "Kansas City"},
	{70, "J", "Ürdün", "Arjantin", utc(2026, 6, 28, 2, 0), "Dallas"},
}

func utc(year int, month time.Month, day, hour, minute int) time.Time {
	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
}

func (m GroupMatch) Deadline() time.Time {
	return m.Kickoff.Add(-1 * time.Hour)
}

func (m GroupMatch) Title() string {
	return SeedTitlePrefix + " Grup " + m.Group + ": " + m.HomeTeam + " vs " + m.AwayTeam
}
