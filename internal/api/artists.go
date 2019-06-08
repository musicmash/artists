package api

import (
	"encoding/json"
	"net/http"

	"github.com/musicmash/artists/internal/db"
	"github.com/musicmash/artists/internal/log"
)

func getArtists(w http.ResponseWriter, r *http.Request) {
	stores, provided := r.URL.Query()["store"]
	if !provided {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(stores[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	artists, err := db.DbMgr.GetArtistsForStore(stores[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(&artists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

// Returns artist_ids that really exist in the db
func validateArtists(w http.ResponseWriter, r *http.Request) {
	artists := []int64{}
	if err := json.NewDecoder(r.Body).Decode(&artists); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(artists) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := db.DbMgr.ValidateArtists(&artists); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(&artists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}
