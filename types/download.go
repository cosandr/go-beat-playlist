package types

const (
	// ScoreSaberStarsURL Scoresaber API URL for getting top X stars
	ScoreSaberStarsURL = "https://scoresaber.com/api.php?function=get-leaderboards&cat=3&limit=%[1]d&page=1&ranked=1"
	// BeatStarAll Dump of all maps
	BeatStarAll = "https://cdn.wes.cloud/beatstar/bssb/v2-all.json"
	// BeatStarRanked Dump of all ranked maps, in desceding PP order
	BeatStarRanked = "https://cdn.wes.cloud/beatstar/bssb/v2-ranked.json"
	// BeatSaverDump Dump of Beatsaver database
	BeatSaverDump = "https://beatsaver.com/api/download/dump/maps"
)

// type ScoreSaberStars struct {
// 	URL string
// 	Raw 
// }
