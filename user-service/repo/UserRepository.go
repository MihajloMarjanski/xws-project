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

	//dsn := "host=userdb user=XML password=ftn dbname=XML_TEST port=5432 sslmode=disable"
	dsn := "host=localhost user=XML password=ftn dbname=XML_TEST port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	repo.db = db
	repo.db.AutoMigrate(&model.User{}, &model.Interest{}, model.Experience{}, model.Block{}, model.JobOffer{})

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

func (repo *UserRepository) SearchUsers(username string) []model.User {
	var users []model.User
	repo.db.Preload("Interests").Preload("Experiences").Model(&users).Where("LOWER(user_name) LIKE ?", "%"+strings.ToLower(username)+"%").Find(&users)
	return users
}

func (repo *UserRepository) GetByID(id int) model.User {
	var user model.User
	repo.db.Preload("Interests").Preload("Experiences").Model(&user).Find(&user, id)
	return user
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

func (repo *UserRepository) CreateUser(name string, email string, password string, username string, gender model.Gender, phonenumber string, dateofbirth time.Time, biography string, apiKey string) int {
	user := model.User{
		Name:        name,
		Email:       email,
		UserName:    username,
		Password:    password,
		Gender:      gender,
		PhoneNumber: phonenumber,
		DateOfBirth: dateofbirth,
		Biography:   biography,
		ApiKey:      apiKey,
	}

	if gender == "Male" || gender == "Female" {
		repo.db.Create(&user)
	} else {
		user.ID = 0
	}
	return int(user.ID)
}

func (repo *UserRepository) UpdateUser(id uint, name string, email string, password string, username string, gender model.Gender, phonenumber string, dateofbirth time.Time, biography string, isPrivate bool) int {
	user := repo.GetByID(int(id))
	user.Name = name
	user.Password = password
	user.Gender = gender
	user.PhoneNumber = phonenumber
	user.DateOfBirth = dateofbirth
	user.Biography = biography
	user.IsPrivate = isPrivate

	if gender == "Male" || gender == "Female" {
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

func (repo *UserRepository) BlockUser(userId int, blockedUserId int) {
	block := model.Block{
		UserId:      uint(userId),
		BlockedUser: uint(blockedUserId)}

	repo.db.Save(&block)
	return
}

func (repo *UserRepository) CreateJobOffer(id int, info string, position string, companyName string, qualifications string, key string) {
	offer := model.JobOffer{
		CompanyName:    companyName,
		JobPosition:    position,
		ApiKey:         key,
		Qualifications: qualifications,
		JobInfo:        info,
	}

	repo.db.Create(&offer)
	return
}

func (repo *UserRepository) ActivateAccount(token string) {
	user := repo.GetByApiKey(token)
	user.IsActivated = true
	repo.db.Save(&user)
}

func (repo *UserRepository) GetByApiKey(token string) model.User {
	var user model.User
	repo.db.Preload("Interests").Preload("Experiences").Model(&user).Where("api_key  = ?", token).Find(&user)
	return user
}

func (repo *UserRepository) FindBlockedForUserId(id uint) []uint {
	var blocked []model.Block
	//repo.db.Raw("SELECT * FROM blocks WHERE user = ?").Scan(&blocked)"%"+strings.ToLower(username)+"%"
	repo.db.Model(&blocked).Where("user_id = ?", id).Find(&blocked)
	var ids []uint
	for _, u := range blocked {
		if !repo.Contains(ids, u.BlockedUser) {
			ids = append(ids, u.BlockedUser)
		}
		if !repo.Contains(ids, u.UserId) {
			ids = append(ids, u.UserId)
		}
	}

	return ids
}

func (repo *UserRepository) Contains(elems []uint, v uint) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
