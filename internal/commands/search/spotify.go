package search

import (
	"sort"

	"github.com/musicmash/artists/internal/db"
	"github.com/musicmash/artists/internal/log"
	"github.com/zmb3/spotify"
)

func parseArtistsAlbums(workerID int, client spotify.Client) {
	log.Infof("worker #%d spawned", workerID)
	for {
		job := <-artistJobs

		log.Infof("worker #%d loading and processing '%s' albums", workerID, job.SpotifyArtist.Name)
		loadAndProcessAlbums(client, job.SpotifyArtist.ID, job.DBArtistID)
		wg.Done()
	}
}

func sortArtistsByPopularity(artists []spotify.FullArtist) []spotify.FullArtist {
	sort.Slice(artists, func(i, j int) bool {
		return artists[i].Popularity > artists[j].Popularity
	})
	return artists
}

func processArtists(artists []spotify.FullArtist) {
	for _, artist := range artists {
		processArtist(artist)
	}
}

func processArtist(artist spotify.FullArtist) {
	if exists := db.DbMgr.IsArtistExistsInStore(storeName, artist.ID.String()); exists && !opts.forceSearchAndSave {
		log.Warn(artist.ID, artist.Name, "already exists")
		return
	}

	newArtist := &db.Artist{
		Name:       artist.Name,
		Popularity: artist.Popularity,
		Followers:  artist.Followers.Count,
	}
	if len(artist.Images) > 0 {
		newArtist.Poster = artist.Images[0].URL
	}

	log.Info("creating new artist", newArtist.Name)
	if err := db.DbMgr.EnsureArtistExists(newArtist); err != nil {
		log.Error("can't create new artist")
	}

	log.Info("saving artist_store_info for new artist", newArtist.ID)
	if err := db.DbMgr.EnsureArtistExistsInStore(newArtist.ID, storeName, artist.ID.String()); err != nil {
		log.Error("can't save spotify id for new artist")
	}

	artistJobs <- &Job{SpotifyArtist: &artist, DBArtistID: newArtist.ID}
	wg.Add(1)
}

func loadAndProcessAlbums(client spotify.Client, artistID spotify.ID, dbArtistID int64) {
	tx := db.DbMgr.Begin()
	albumPage, err := client.GetArtistAlbums(artistID)
	if err != nil {
		log.Error(err)
		tx.Rollback()
	}

	processAlbums(albumPage.Albums, dbArtistID, tx)

	for albumPage.Total > albumPage.Limit+albumPage.Offset {
		albumPage.Offset += albumPage.Limit
		log.Infof("getting next albums for artist %v...", dbArtistID)

		opts := spotify.Options{
			Limit:  &limit,
			Offset: &albumPage.Offset,
		}
		albumPage, err = client.GetArtistAlbumsOpt(artistID, &opts, nil)
		if err != nil {
			tx.Commit()
			log.Panic(err)
		}

		processAlbums(albumPage.Albums, dbArtistID, tx)
	}
	tx.Commit()
}

func processAlbums(albums []spotify.SimpleAlbum, dbArtistID int64, tx db.DataMgr) {
	for _, album := range albums {
		log.Debugf("process albums from %s", album.Artists[0].Name)
		processAlbum(album, dbArtistID, tx)
	}
}

func processAlbum(album spotify.SimpleAlbum, dbArtistID int64, tx db.DataMgr) {
	log.Debugf("saving album %s", album.Name)
	err := tx.EnsureAlbumExists(&db.Album{
		ArtistID: dbArtistID,
		Name:     album.Name,
	})

	if err != nil {
		log.Error(err)
	}

	// handle other artists mentioned in this album
	//for _, artist := range album.Artists {
	//	processArtist(client, spotify.FullArtist{SimpleArtist:artist})
	//}
}
