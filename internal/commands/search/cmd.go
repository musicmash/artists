package search

import (
	"context"
	"errors"
	"sync"

	"github.com/musicmash/artists/internal/db"
	"github.com/musicmash/artists/internal/log"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

const storeName = "spotify"

var (
	artistJobs chan *Job
	wg         = sync.WaitGroup{}
	limit      = 50
)

type Job struct {
	SpotifyArtist *spotify.FullArtist
	DBArtistID    int64
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search artists on Spotify and save their albums on the db",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("artist name for search not provided")
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			db.DbMgr = db.NewMainDatabaseMgr()
			artistJobs = make(chan *Job, opts.workersCount)

			log.Info("ensuring that 'spotify' exists...")
			return db.DbMgr.EnsureStoreExists(storeName)
		},
		RunE: run,
		PostRun: func(cmd *cobra.Command, args []string) {
			close(artistJobs)
		},
	}

	cmd.Flags().StringVar(&opts.clientID, "spotify-app-id", "", "spotify app id")
	_ = cmd.MarkFlagRequired("spotify-app-id")
	cmd.Flags().StringVar(&opts.clientSecret, "spotify-app-secret", "", "spotify app secret")
	_ = cmd.MarkFlagRequired("spotify-app-secret")
	cmd.Flags().IntVar(&opts.workersCount, "workers", 5, "count of workers")
	cmd.Flags().BoolVar(&opts.forceSearchAndSave, "force", false, "search artist and overwrite if exists")
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	credentials := &clientcredentials.Config{
		ClientID:     opts.clientID,
		ClientSecret: opts.clientSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := credentials.Token(context.Background())
	if err != nil {
		return err
	}

	client := spotify.Authenticator{}.NewClient(token)
	for workerID := 1; workerID <= opts.workersCount; workerID++ {
		go parseArtistsAlbums(workerID, client)
	}

	log.Infof("searching '%s'", args[0])
	results, err := client.SearchOpt(args[0], spotify.SearchTypeArtist, &spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return err
	}
	log.Debugf("limit %v offset %v total %v", results.Artists.Limit, results.Artists.Offset, results.Artists.Total)
	processArtists(client, sortArtistsByPopularity(results.Artists.Artists))

	// load next part
	for results.Artists.Total > results.Artists.Limit+results.Artists.Offset {
		log.Info("getting next artists...")
		if err = client.NextArtistResults(results); err != nil {
			return err
		}

		log.Debugf("limit %v offset %v total %v", results.Artists.Limit, results.Artists.Offset, results.Artists.Total)
		processArtists(client, sortArtistsByPopularity(results.Artists.Artists))
	}

	wg.Wait()
	return nil
}
