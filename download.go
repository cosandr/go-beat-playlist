package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	// The user agent used for HTTP GET requests
	httpUserAgent = "go_beat_playlist/1.0"
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

var httpClient = &http.Client{}

func httpGet(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", httpUserAgent)
	resp, err = httpClient.Do(req)
	return
}

// DownloadSong tries to download a song from BeatSaver using its hash or key, returns a DownloadSong
//
// Function merges downloaded metadata with argument, downloaded song is saved to `path`
func DownloadSong(s *Song) (retSong Song, err error) {
	// Working Song
	var dlSong Song
	if len(s.URL) == 0 {
		bsSong, errDl := DownloadSongInfo(s)
		if errDl != nil {
			err = errDl
			return
		}
		dlSong = bsSong
	} else {
		dlSong = *s
	}
	dlPath := fmt.Sprintf("%s/%s", conf.Songs, dlSong.DirName())
	if !DirExists(dlPath) {
		songBytes, errB := DownloadSongBytes(dlSong.URL)
		if errB != nil {
			err = errB
			return
		}
		errB = ExtractZIP(dlPath, &songBytes)
		if errB != nil {
			err = errB
			return
		}
	}
	// Load downloaded song
	infoPath, err := FindInfo(dlPath)
	if err != nil {
		return
	}
	retSong, err = MakeSong(infoPath)
	if err != nil {
		return
	}
	if dlSong.Hash != retSong.Hash {
		err = fmt.Errorf("download failed, hash mismatch")
		os.RemoveAll(dlPath)
		return
	}
	retSong = retSong.Merge(&dlSong)
	return
}

// DownloadSongBytes tries to download a song from BeatSaver using its url, returns byte array
func DownloadSongBytes(url string) (out []byte, err error) {
	dl, err := httpGet("https://beatsaver.com" + url)
	if err != nil {
		return
	}
	defer dl.Body.Close()
	out, err = ioutil.ReadAll(dl.Body)
	return
}

// ExtractZIP extract byte slice (ZIP file) to `path`
func ExtractZIP(path string, in *[]byte) (err error) {
	if !DirExists(path) {
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

// DownloadSongInfo fetches song info from BeatSaver API, returns a new Song
func DownloadSongInfo(s *Song) (dlSong Song, err error) {
	var url string
	if len(s.Hash) > 0 {
		url = fmt.Sprintf(beatSaverByHash, s.Hash)
	} else if len(s.Key) > 0 {
		url = fmt.Sprintf(beatSaverByKey, s.Key)
	} else {
		err = fmt.Errorf("%s has no key or hash", s.Name)
		return
	}
	resp, err := httpGet(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("HTTP GET failed: %s", resp.Status)
		return
	}
	outSong, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	dlSong, err = MakeBeatSaverSong(&outSong)
	if err != nil {
		return
	}
	return
}

// DownloadStarsPlaylist returns a Playlist of top `num` songs sorted by star difficulty
func DownloadStarsPlaylist(num int) (p Playlist, err error) {
	resp, err := httpGet(fmt.Sprintf(scoreSaberStarsURL, num))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	p, err = MakeScoreSaberPlaylist(&body)
	if err != nil {
		return
	}
	if len(p.Songs) == 0 {
		err = fmt.Errorf("response parsing failed")
		return
	}
	return
}

// DownloadPPPlaylist returns a Playlist of top `num` songs sorted by PP
func DownloadPPPlaylist(num int) (p Playlist, err error) {
	resp, err := httpGet(beatStarRanked)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	p, err = MakeSongBrowserPlaylist(&body)
	if err != nil {
		return
	}
	// Sort by PP
	p.SortByPP()
	// Only keep num songs
	p.Songs = p.Songs[:num]
	return
}
