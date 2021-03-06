package db

import "github.com/jinzhu/gorm"

var tables = []interface{}{
	Artist{},
	ArtistStoreInfo{},
	Store{},
	Album{},
}

func CreateTables(db *gorm.DB) error {
	return db.AutoMigrate(tables...).Error
}

func DropAllTables(db *gorm.DB) error {
	return db.DropTable(tables...).Error
}

func CreateAll(db *gorm.DB) error {
	if err := CreateTables(db); err != nil {
		return err
	}

	fkeys := map[interface{}][][2]string{
		&ArtistStoreInfo{}: {
			{"artist_id", "artists(id)"},
			{"store_name", "stores(name)"},
		},
		&Album{}: {
			{"artist_id", "artists(id)"},
		},
	}

	for model, foreignKey := range fkeys {
		for _, fk := range foreignKey {
			if err := db.Debug().Model(model).AddForeignKey(
				fk[0], fk[1], "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
		}
	}

	if err := db.Debug().Model(&Album{}).AddUniqueIndex(
		"idx_album_art_id_name",
		"artist_id", "name").Error; err != nil {
		return err
	}

	if err := db.Debug().Model(&ArtistStoreInfo{}).AddUniqueIndex(
		"idx_art_store_name_id",
		"artist_id", "store_name", "store_id").Error; err != nil {
		return err
	}

	return nil
}
