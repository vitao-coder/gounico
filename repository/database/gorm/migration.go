package gorm

import "gounico/feiralivre/domain"

func (db *Database) AutoMigrateDatabase() error {

	err := db.autoMigrate(&domain.RegiaoGenerica{})
	if err != nil {
		return err
	}

	err = db.autoMigrate(&domain.Regiao{})
	if err != nil {
		return err
	}

	err = db.autoMigrate(&domain.Distrito{})
	if err != nil {
		return err
	}

	err = db.autoMigrate(&domain.Localizacao{})
	if err != nil {
		return err
	}

	err = db.autoMigrate(&domain.SubPrefeitura{})
	if err != nil {
		return err
	}

	err = db.autoMigrate(&domain.Feira{})
	if err != nil {
		return err
	}
	return nil
}
