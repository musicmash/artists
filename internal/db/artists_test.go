package db

import (
	"testing"

	"github.com/musicmash/artists/internal/testutil/vars"
	"github.com/stretchr/testify/assert"
)

func TestDB_Artist_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex})

	// assert
	assert.NoError(t, err)
}

func TestDB_Artists_GetAll(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))

	// action
	artists, err := DbMgr.GetAllArtists()

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}

func TestDB_Artists_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex, Followers: 100}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistArchitects, Followers: 250}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistSPY}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistWildways, Followers: 50}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{Name: vars.ArtistRitaOra, Followers: 90}))
	want := []struct {
		SearchText string
		Artists    []string
	}{
		{SearchText: "il", Artists: []string{vars.ArtistSkrillex, vars.ArtistWildways}},
		{SearchText: vars.ArtistSkrillex, Artists: []string{vars.ArtistSkrillex}},
		{SearchText: "a", Artists: []string{vars.ArtistArchitects, vars.ArtistRitaOra, vars.ArtistWildways}},
	}

	for _, item := range want {
		// action
		artists, err := DbMgr.SearchArtists(item.SearchText)

		// assert
		assert.NoError(t, err)
		assert.Len(t, artists, len(item.Artists))
		for i, wantName := range item.Artists {
			assert.Equal(t, wantName, artists[i].Name)
		}
	}
}

func TestDB_ArtistStoreInfo_EnsureArtistExistsInStore(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureStoreExists(vars.StoreDeezer))

	// action
	err := DbMgr.EnsureArtistExistsInStore(vars.StoreIDQ, vars.StoreDeezer, vars.StoreIDA)

	// assert
	assert.NoError(t, err)
	artists, err := DbMgr.GetArtistsForStore(vars.StoreDeezer)
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}

func TestDB_ArtistStoreInfo_GetArtistFromStore(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureStoreExists(vars.StoreApple))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(vars.StoreIDQ, vars.StoreApple, vars.StoreIDA))
	assert.NoError(t, DbMgr.EnsureArtistExistsInStore(vars.StoreIDQ, vars.StoreApple, vars.StoreIDB))

	// action
	artists, err := DbMgr.GetArtistFromStore(vars.StoreIDQ, vars.StoreApple)

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
}

func TestDB_Artists_Validate(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	artists := []int64{1, 5, 10, 420, 69, 333, 999}
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 1, Name: vars.ArtistSkrillex, Followers: 100}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 2, Name: vars.ArtistArchitects, Followers: 250}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 3, Name: vars.ArtistSPY}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 4, Name: vars.ArtistWildways, Followers: 50}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 5, Name: vars.ArtistRitaOra, Followers: 90}))

	// action
	err := DbMgr.ValidateArtists(&artists)

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0])
	assert.Equal(t, int64(5), artists[1])
}

func TestDB_Artists_GetWithFullInfo(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 1, Name: vars.ArtistSkrillex, Followers: 100}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 2, Name: vars.ArtistArchitects, Followers: 250}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 3, Name: vars.ArtistSPY}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 4, Name: vars.ArtistWildways, Followers: 50}))
	assert.NoError(t, DbMgr.EnsureArtistExists(&Artist{ID: 5, Name: vars.ArtistRitaOra, Followers: 90}))

	// action
	artists, err := DbMgr.GetArtistsWithFullInfo([]int64{1, 5})

	// assert
	assert.NoError(t, err)
	assert.Len(t, artists, 2)

	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, vars.ArtistSkrillex, artists[0].Name)
	assert.Equal(t, uint(100), artists[0].Followers)

	assert.Equal(t, int64(5), artists[1].ID)
	assert.Equal(t, vars.ArtistRitaOra, artists[1].Name)
	assert.Equal(t, uint(90), artists[1].Followers)
}
