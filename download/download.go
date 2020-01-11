package download

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	// beatSaverByKey URL to download from key
	beatSaverByKey = "https://beatsaver.com/api/maps/detail/%s"
	// beatSaverByHash URL to download from hash
	beatSaverByHash = "https://beatsaver.com/api/maps/by-hash/%s"
)

// Song tries to download a song from BeatSaver using its hash or key, returns a Song
//
// Function merges downloaded metadata with argument, downloaded song is saved to `path`
func Song(path string, s *mt.Song) (retSong mt.Song, err error) {
	// Working Song
	var dlSong mt.Song
	if len(s.URL) == 0 {
		bsSong, errDl := SongInfo(s)
		if errDl != nil {
			err = errDl
			return
		}
		dlSong = bsSong
	} else {
		dlSong = *s
	}
	if !mt.DirExists(path) {
		songBytes, errB := SongBytes(dlSong.URL)
		if errB != nil {
			err = errB
			return
		}
		errB = ExtractZIP(path, &songBytes)
		if errB != nil {
			err = errB
			return
		}
	}
	// Load downloaded song
	retSong, err = mt.MakeSong(path + "/info.dat")
	if err != nil {
		return
	}
	if dlSong.Hash != retSong.Hash {
		err = fmt.Errorf("download failed, hash mismatch")
		return
	}
	retSong = retSong.Merge(&dlSong)
	return
}

// SongBytes tries to download a song from BeatSaver using its url, returns byte array
func SongBytes(url string) (out []byte, err error) {
	dl, err := http.Get("https://beatsaver.com" + url)
	if err != nil {
		return
	}
	defer dl.Body.Close()
	out, err = ioutil.ReadAll(dl.Body)
	return
}

// ExtractZIP extract byte slice (ZIP file) to `path`
func ExtractZIP(path string, in *[]byte) (err error) {
	if !mt.DirExists(path) {
		errMk := os.MkdirAll(path, 0755)
		if errMk != nil {
			err = errMk
			return
		}
	}
	zipReader, err := zip.NewReader(bytes.NewReader(*in), int64(len(*in)))
	if err != nil {
		return
	}
	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = ioutil.WriteFile(path+"/"+zipFile.Name, unzippedFileBytes, 0755)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// SongInfo fetches song info from BeatSaver API, returns a new Song
func SongInfo(s *mt.Song) (dlSong mt.Song, err error) {
	var url string
	if len(s.Hash) > 0 {
		url = fmt.Sprintf(beatSaverByHash, s.Hash)
	} else if len(s.Key) > 0 {
		url = fmt.Sprintf(beatSaverByKey, s.Key)
	} else {
		err = fmt.Errorf("%s has no key or hash", s.Name)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	outSong, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	dlSong, err = mt.MakeBeatSaverSong(&outSong)
	if err != nil {
		return
	}
	return
}

// StarsPlaylist returns a Playlist of top `num` songs sorted by star difficulty
func StarsPlaylist(num int) (p mt.Playlist, err error) {
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

// PPPlaylist returns a Playlist of top `num` songs sorted by PP
func PPPlaylist(num int) (p mt.Playlist, err error) {
	resp, err := http.Get(beatStarRanked)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	p, err = mt.MakeSongBrowserPlaylist(&body)
	if err != nil {
		return
	}
	// Sort by PP
	p.SortByPP()
	// Only keep num songs
	p.Songs = p.Songs[:num]
	return
}
