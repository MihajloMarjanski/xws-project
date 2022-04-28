package repo

import (
	"time"
	"user-service/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func New() (*UserRepository, error) {
	repo := &UserRepository{}

	dsn := "host=localhost user=XML password=ftn dbname=XML_TEST port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	repo.db = db
	repo.db.AutoMigrate(&model.User{})

	return repo, nil
}

func (repo *UserRepository) Close() error {
	db, err := repo.db.DB()
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func (repo *UserRepository) CreateUser(name string, email string, password string, username string, gender model.Gender, phonenumber string, dateofbirth time.Time, biography string) int {
	user := model.User{
		Name:        name,
		Email:       email,
		UserName:    username,
		Password:    password,
		Gender:      gender,
		PhoneNumber: phonenumber,
		DateOfBirth: dateofbirth,
		Biography:   biography}

	if gender == "male" || gender == "female" {
		repo.db.Create(&user)
	} else {
		user.ID = 0
	}
	return int(user.ID)
}
