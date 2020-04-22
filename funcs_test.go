package main

import "testing"

func TestMakePlaylist(t *testing.T) {
	p, err := MakePlaylist("samples/json/playlist.bplist")
	if err != nil {
		t.Errorf("Playlist JSON parse failed: %v", err)
	}
	if len(p.Songs) != 44 {
		t.Errorf("Expected 44 songs, got %d\n%s", len(p.Songs), p.Debug())
	}
}
