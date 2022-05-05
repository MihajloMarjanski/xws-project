package repo

import (
	"strings"
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
	repo.db.AutoMigrate(&model.User{}, &model.Interest{}, model.Experience{})

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

func (repo *UserRepository) SearchUsers(username string) []model.UserDTO {
	var users []model.User
	var usersDTO []model.UserDTO
	repo.db.Preload("Interests").Preload("Experiences").Model(&users).Where("LOWER(user_name) LIKE ?", "%"+strings.ToLower(username)+"%").Find(&usersDTO)
	return usersDTO
}

func (repo *UserRepository) GetByID(id int) model.UserDTO {
	var user model.User
	var userDTO model.UserDTO
	repo.db.Preload("Interests").Preload("Experiences").Model(&user).Find(&userDTO, id)
	return userDTO
}

func (repo *UserRepository) GetByUsername(username string) model.User {
	var user model.User
	repo.db.Preload("Interests").Preload("Experiences").Model(&user).Where("user_name  = ?", username).Find(&user)
	return user
}

func (repo *UserRepository) GetMe(id int) model.User {
	var user model.User
	repo.db.Preload("Interests").Preload("Experiences").Model(&user).Find(&user, id)
	return user
}

func (repo *UserRepository) RemoveExperience(id int) int {
	repo.db.Delete(&model.Experience{}, id)
	return id
}

func (repo *UserRepository) RemoveInterest(id int) int {
	repo.db.Delete(&model.Interest{}, id)
	return id
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

func (repo *UserRepository) UpdateUser(id uint, name string, email string, password string, username string, gender model.Gender, phonenumber string, dateofbirth time.Time, biography string) int {
	user := model.User{
		ID:          uint(id),
		Name:        name,
		Email:       email,
		UserName:    username,
		Password:    password,
		Gender:      gender,
		PhoneNumber: phonenumber,
		DateOfBirth: dateofbirth,
		Biography:   biography}

	if gender == "male" || gender == "female" {
		repo.db.Save(&user)
	} else {
		user.ID = 0
	}
	return int(user.ID)
}

func (repo *UserRepository) AddInterest(interest string, userId uint) int {
	newInterest := model.Interest{
		Interest: interest,
		UserID:   userId}

	repo.db.Create(&newInterest)

	return int(newInterest.ID)
}

func (repo *UserRepository) AddExperience(company string, position string, from time.Time, until time.Time, userId uint) int {
	experience := model.Experience{
		Company:  company,
		Position: position,
		From:     from,
		Until:    until,
		UserID:   userId}

	repo.db.Create(&experience)

	return int(experience.ID)
}
