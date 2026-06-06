package worldcup2026

import "time"

type GroupMatch struct {
	MatchNo   int
	Group     string
	HomeTeam  string
	AwayTeam  string
	Kickoff   time.Time
	Venue     string
}

// GroupMatches contains all 72 group-stage fixtures.
// Kickoff times are UTC (GMT). Deadlines for predictions are kickoff minus 1 hour.
var GroupMatches = []GroupMatch{
	{1, "A", "Mexico", "South Africa", utc(2026, 6, 11, 19, 0), "Mexico City"},
	{2, "A", "Korea Republic", "Czechia", utc(2026, 6, 12, 2, 0), "Guadalajara"},
	{3, "B", "Canada", "Bosnia and Herzegovina", utc(2026, 6, 12, 19, 0), "Toronto"},
	{4, "D", "USA", "Paraguay", utc(2026, 6, 13, 1, 0), "Los Angeles"},
	{8, "B", "Qatar", "Switzerland", utc(2026, 6, 13, 19, 0), "San Francisco Bay Area"},
	{7, "C", "Brazil", "Morocco", utc(2026, 6, 13, 22, 0), "New York/New Jersey"},
	{5, "C", "Haiti", "Scotland", utc(2026, 6, 14, 1, 0), "Boston"},
	{6, "D", "Australia", "Türkiye", utc(2026, 6, 14, 4, 0), "Vancouver"},
	{10, "E", "Germany", "Curaçao", utc(2026, 6, 14, 17, 0), "Houston"},
	{11, "F", "Netherlands", "Japan", utc(2026, 6, 14, 20, 0), "Dallas"},
	{9, "E", "Côte d'Ivoire", "Ecuador", utc(2026, 6, 14, 23, 0), "Philadelphia"},
	{12, "F", "Sweden", "Tunisia", utc(2026, 6, 15, 2, 0), "Monterrey"},
	{14, "H", "Spain", "Cabo Verde", utc(2026, 6, 15, 16, 0), "Atlanta"},
	{16, "G", "Belgium", "Egypt", utc(2026, 6, 15, 19, 0), "Seattle"},
	{13, "H", "Saudi Arabia", "Uruguay", utc(2026, 6, 15, 22, 0), "Miami"},
	{15, "G", "IR Iran", "New Zealand", utc(2026, 6, 16, 1, 0), "Los Angeles"},
	{17, "I", "France", "Senegal", utc(2026, 6, 16, 19, 0), "New York/New Jersey"},
	{18, "I", "Iraq", "Norway", utc(2026, 6, 16, 22, 0), "Boston"},
	{19, "J", "Argentina", "Algeria", utc(2026, 6, 17, 1, 0), "Kansas City"},
	{20, "J", "Austria", "Jordan", utc(2026, 6, 17, 4, 0), "San Francisco Bay Area"},
	{23, "K", "Portugal", "Congo DR", utc(2026, 6, 17, 17, 0), "Houston"},
	{22, "L", "England", "Croatia", utc(2026, 6, 17, 20, 0), "Dallas"},
	{21, "L", "Ghana", "Panama", utc(2026, 6, 17, 23, 0), "Toronto"},
	{24, "K", "Uzbekistan", "Colombia", utc(2026, 6, 18, 2, 0), "Mexico City"},
	{25, "A", "Czechia", "South Africa", utc(2026, 6, 18, 16, 0), "Atlanta"},
	{26, "B", "Switzerland", "Bosnia and Herzegovina", utc(2026, 6, 18, 19, 0), "Los Angeles"},
	{27, "B", "Canada", "Qatar", utc(2026, 6, 18, 22, 0), "Vancouver"},
	{28, "A", "Mexico", "Korea Republic", utc(2026, 6, 19, 1, 0), "Guadalajara"},
	{32, "D", "USA", "Australia", utc(2026, 6, 19, 19, 0), "Seattle"},
	{30, "C", "Scotland", "Morocco", utc(2026, 6, 19, 22, 0), "Boston"},
	{29, "C", "Brazil", "Haiti", utc(2026, 6, 20, 0, 30), "Philadelphia"},
	{31, "D", "Türkiye", "Paraguay", utc(2026, 6, 20, 3, 0), "San Francisco Bay Area"},
	{35, "F", "Netherlands", "Sweden", utc(2026, 6, 20, 17, 0), "Houston"},
	{33, "E", "Germany", "Côte d'Ivoire", utc(2026, 6, 20, 20, 0), "Toronto"},
	{34, "E", "Ecuador", "Curaçao", utc(2026, 6, 21, 0, 0), "Kansas City"},
	{36, "F", "Tunisia", "Japan", utc(2026, 6, 21, 4, 0), "Monterrey"},
	{38, "H", "Spain", "Saudi Arabia", utc(2026, 6, 21, 16, 0), "Atlanta"},
	{39, "G", "Belgium", "IR Iran", utc(2026, 6, 21, 19, 0), "Los Angeles"},
	{37, "H", "Uruguay", "Cabo Verde", utc(2026, 6, 21, 22, 0), "Miami"},
	{40, "G", "New Zealand", "Egypt", utc(2026, 6, 22, 1, 0), "Vancouver"},
	{43, "J", "Argentina", "Austria", utc(2026, 6, 22, 17, 0), "Dallas"},
	{42, "I", "France", "Iraq", utc(2026, 6, 22, 21, 0), "Philadelphia"},
	{41, "I", "Norway", "Senegal", utc(2026, 6, 23, 0, 0), "New York/New Jersey"},
	{44, "J", "Jordan", "Algeria", utc(2026, 6, 23, 3, 0), "San Francisco Bay Area"},
	{47, "K", "Portugal", "Uzbekistan", utc(2026, 6, 23, 17, 0), "Houston"},
	{45, "L", "England", "Ghana", utc(2026, 6, 23, 20, 0), "Boston"},
	{46, "L", "Panama", "Croatia", utc(2026, 6, 23, 23, 0), "Toronto"},
	{48, "K", "Colombia", "Congo DR", utc(2026, 6, 24, 2, 0), "Guadalajara"},
	{51, "B", "Switzerland", "Canada", utc(2026, 6, 24, 19, 0), "Vancouver"},
	{52, "B", "Bosnia and Herzegovina", "Qatar", utc(2026, 6, 24, 19, 0), "Seattle"},
	{49, "C", "Scotland", "Brazil", utc(2026, 6, 24, 22, 0), "Miami"},
	{50, "C", "Morocco", "Haiti", utc(2026, 6, 24, 22, 0), "Atlanta"},
	{53, "A", "Czechia", "Mexico", utc(2026, 6, 25, 1, 0), "Mexico City"},
	{54, "A", "South Africa", "Korea Republic", utc(2026, 6, 25, 1, 0), "Monterrey"},
	{55, "E", "Curaçao", "Côte d'Ivoire", utc(2026, 6, 25, 20, 0), "Philadelphia"},
	{56, "E", "Ecuador", "Germany", utc(2026, 6, 25, 20, 0), "New York/New Jersey"},
	{57, "F", "Japan", "Sweden", utc(2026, 6, 25, 23, 0), "Dallas"},
	{58, "F", "Tunisia", "Netherlands", utc(2026, 6, 25, 23, 0), "Kansas City"},
	{59, "D", "Türkiye", "USA", utc(2026, 6, 26, 2, 0), "Los Angeles"},
	{60, "D", "Paraguay", "Australia", utc(2026, 6, 26, 2, 0), "San Francisco Bay Area"},
	{61, "I", "Norway", "France", utc(2026, 6, 26, 19, 0), "Boston"},
	{62, "I", "Senegal", "Iraq", utc(2026, 6, 26, 19, 0), "Toronto"},
	{65, "H", "Cabo Verde", "Saudi Arabia", utc(2026, 6, 27, 0, 0), "Houston"},
	{66, "H", "Uruguay", "Spain", utc(2026, 6, 27, 0, 0), "Guadalajara"},
	{63, "G", "Egypt", "IR Iran", utc(2026, 6, 27, 3, 0), "Seattle"},
	{64, "G", "New Zealand", "Belgium", utc(2026, 6, 27, 3, 0), "Vancouver"},
	{67, "L", "Panama", "England", utc(2026, 6, 27, 21, 0), "New York/New Jersey"},
	{68, "L", "Croatia", "Ghana", utc(2026, 6, 27, 21, 0), "Philadelphia"},
	{71, "K", "Colombia", "Portugal", utc(2026, 6, 27, 23, 30), "Miami"},
	{72, "K", "Congo DR", "Uzbekistan", utc(2026, 6, 27, 23, 30), "Atlanta"},
	{69, "J", "Algeria", "Austria", utc(2026, 6, 28, 2, 0), "Kansas City"},
	{70, "J", "Jordan", "Argentina", utc(2026, 6, 28, 2, 0), "Dallas"},
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
