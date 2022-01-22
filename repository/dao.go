package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daniarmas/api-example/models"
	"github.com/daniarmas/api-example/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DAO interface {
	NewItemQuery() ItemQuery
	NewUserQuery() UserQuery
	NewRefreshTokenQuery() RefreshTokenQuery
	NewAuthorizationTokenQuery() AuthorizationTokenQuery
	NewTokenQuery() TokenQuery
	NewHashPasswordQuery() HashPasswordQuery
}

type dao struct{}

var DB *gorm.DB
var Config *utils.Config

func NewDAO(db *gorm.DB, config *utils.Config) DAO {
	DB = db
	Config = config
	return &dao{}
}

func NewConfig() (*utils.Config, error) {
	Config, err := utils.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	return &Config, nil
}

func NewDB(config *utils.Config) (*gorm.DB, error) {
	host := config.DBHost
	port := config.DBPort
	user := config.DBUser
	dbName := config.DBDatabase
	password := config.DBPassword

	// Starting a database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 200,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	DB.AutoMigrate(&models.User{}, &models.Item{}, &models.RefreshToken{}, &models.AuthorizationToken{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func (d *dao) NewItemQuery() ItemQuery {
	return &itemQuery{}
}

func (d *dao) NewUserQuery() UserQuery {
	return &userQuery{}
}

func (d *dao) NewRefreshTokenQuery() RefreshTokenQuery {
	return &refreshTokenQuery{}
}

func (d *dao) NewAuthorizationTokenQuery() AuthorizationTokenQuery {
	return &authorizationTokenQuery{}
}

func (d *dao) NewTokenQuery() TokenQuery {
	return &tokenQuery{}
}

func (d *dao) NewHashPasswordQuery() HashPasswordQuery {
	return &hashPasswordQuery{}
}
