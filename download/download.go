package download

import (
	"fmt"
	"io/ioutil"
	"net/http"

	mt "github.com/cosandr/go-beat-playlist/types"
)

const (
	// ScoreSaberStarsURL Scoresaber API URL for getting top X stars
	scoreSaberStarsURL = "https://scoresaber.com/api.php?function=get-leaderboards&cat=3&limit=%[1]d&page=1&ranked=1"
	// BeatStarAll Dump of all maps
	beatStarAll = "https://cdn.wes.cloud/beatstar/bssb/v2-all.json"
	// BeatStarRanked Dump of all ranked maps, in desceding PP order
	beatStarRanked = "https://cdn.wes.cloud/beatstar/bssb/v2-ranked.json"
	// BeatSaverDump Dump of Beatsaver database
	beatSaverDump = "https://beatsaver.com/api/download/dump/maps"
)

// FetchByStars returns a Playlist of top `num` songs sorted by star difficulty
func FetchByStars(num int) (p mt.Playlist, err error) {
	resp, err := http.Get(fmt.Sprintf(scoreSaberStarsURL, num))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	p, err = mt.MakeScoreSaberPlaylist(&body)
	if err != nil {
		return
	}
	if len(p.Songs) == 0 {
		err = fmt.Errorf("response parsing failed")
		return
	}
	return
}
