package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/artists/pkg/api"
)

func Do(provider *api.Provider, artistName string) ([]*Artist, error) {
	searchURL := fmt.Sprintf("%s/search?artist_name=%s", provider.URL, url.QueryEscape(artistName))
	resp, err := provider.Client.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("got %d status code", resp.StatusCode)
	}

	artists := []*Artist{}
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}
