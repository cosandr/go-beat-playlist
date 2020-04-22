package main

// ConfigJSON is the structure of the config.json file
type ConfigJSON struct {
	Game string `json:"game"`
}

// PlaylistJSON is the structure of a playlist JSON or BPLIST
type PlaylistJSON struct {
	Title       string     `json:"playlistTitle"`
	Author      string     `json:"playlistAuthor"`
	Description string     `json:"playlistDescription,omitempty"`
	Image       string     `json:"image,omitempty"`
	Count       int        `json:"playlistSongCount,omitempty"`
	Songs       []SongJSON `json:"songs"`
}

// SongJSON is the structure of a song in a playlist JSON
type SongJSON struct {
	Key      interface{} `json:"key,omitempty"`
	Hash     string      `json:"hash"`
	Name     string      `json:"songName"`
	Uploader string      `json:"uploader,omitempty"`
}

// InfoJSON is the structure of a song's info.dat file (only relevant bits)
type InfoJSON struct {
	SongName   string           `json:"_songName"`
	SongAuthor string           `json:"_songAuthorName"`
	Mapper     string           `json:"_levelAuthorName"`
	Beatmaps   []BeatmapSetJSON `json:"_difficultyBeatmapSets"`
}

// BeatmapSetJSON is the type used to store beatmap sets (types)
type BeatmapSetJSON struct {
	Type string        `json:"_beatmapCharacteristicName"`
	Maps []BeatmapJSON `json:"_difficultyBeatmaps"`
}

// BeatmapJSON is the type for a mapping of a difficulty
type BeatmapJSON struct {
	Difficulty string `json:"_difficulty"`
	File       string `json:"_beatmapFilename"`
}
