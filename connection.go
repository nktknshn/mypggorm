package mypggorm

import (
	"gorm.io/gorm"
)

type DatabaseConnection struct {
	connection *gorm.DB
}

func NewDatabaseConnection(connection *gorm.DB) *DatabaseConnection {
	return &DatabaseConnection{connection: connection}
}

func (dc *DatabaseConnection) Connection() *gorm.DB {
	return dc.connection
}

func (dc *DatabaseConnection) JustExec(sql string) error {
	return dc.connection.Exec(sql).Error
}

func (dc *DatabaseConnection) Close() error {
	sqlDB, err := dc.connection.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
