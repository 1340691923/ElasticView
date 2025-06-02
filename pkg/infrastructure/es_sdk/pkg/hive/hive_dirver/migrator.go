package hive_dirver

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

type Migrator struct {
	migrator.Migrator
	*Dialector
}

func (m Migrator) CurrentDatabase() (name string) {
	m.handleError(m.DB.Raw("SELECT current_database()").Row().Scan(&name))
	return
}

func (m Migrator) HasTable(value interface{}) bool {
	var name string
	m.handleError(m.RunWithValue(value, func(stmt *gorm.Statement) error {
		//currentDatabase := m.DB.Migrator().CurrentDatabase()
		return m.DB.Raw(
			"SHOW TABLES LIKE ?", stmt.Table,
		).Row().Scan(&name)
	}), sql.ErrNoRows)
	return name != ""
}
