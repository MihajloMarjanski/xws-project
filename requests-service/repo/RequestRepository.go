package repo

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"requests-service/model"
	"time"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RequestsRepository struct {
	db *gorm.DB
}

type IsPublicModel struct {
	Id uint
}

func New() (*RequestsRepository, error) {

	repo := &RequestsRepository{}

	dsn := "host=requestdb user=XML password=ftn dbname=XML_REQUESTS port=5432 sslmode=disable"
	//dsn := "host=localhost user=XML password=ftn dbname=XML_REQUESTS port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	repo.db = db
	repo.db.AutoMigrate(&model.Request{}, &model.Connection{}, &model.Message{}, &model.Notification{})

	return repo, nil
}

func (repo *RequestsRepository) Close() error {
	db, err := repo.db.DB()
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func (repo *RequestsRepository) GetAllByRecieverId(rid uint) []model.Request {
	var request model.Request
	var requests []model.Request
	repo.db.Model(&request).Where("receiver_id  = ?", rid).Find(&requests)
	return requests
}

func (repo *RequestsRepository) AcceptRequest(sid, rid uint) {
	request := model.Request{
		SenderID:   sid,
		ReceiverID: rid,
	}
	repo.db.Delete(&request)

	connection := model.Connection{
		UserOne: sid,
		UserTwo: rid,
	}

	repo.db.Create(&connection)
}

func (repo *RequestsRepository) DeclineRequest(sid, rid uint) {
	request := model.Request{
		SenderID:   sid,
		ReceiverID: rid,
	}
	repo.db.Delete(&request)
}

func (repo *RequestsRepository) SendRequest(sid, rid uint) {
	conn, err := grpc.Dial("user-service:8100", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	response, err := client.GetPrivateStatusForUserId(context.Background(), &pb.GetPrivateStatusForUserIdRequest{Id: int64(rid)})
	if err != nil {
		panic(err)
	}

	isPrivate := response.IsPrivate
	if isPrivate {
		request := model.Request{
			SenderID:   sid,
			ReceiverID: rid,
		}

		repo.db.Create(&request)
	} else {
		connection := model.Connection{
			UserOne: sid,
			UserTwo: rid,
		}

		repo.db.Create(&connection)

	}
}

func (repo *RequestsRepository) SendMessage(senderID, receiverID uint, text string) {
	message := model.Message{
		Text:       text,
		SenderId:   senderID,
		ReceiverId: receiverID,
	}
	repo.db.Create(&message)
}

func (repo *RequestsRepository) AreConnected(senderID, receiverID uint) bool {
	var connection model.Connection
	repo.db.Model(&connection).Where("user_one = ? AND user_two = ? OR user_one = ? AND user_two = ?", senderID, receiverID, receiverID, senderID).Find(&connection)
	if connection.UserTwo == 0 {
		return false
	}
	return true
}

func (repo *RequestsRepository) AreRequested(u1 uint, u2 uint) bool {
	var request model.Request
	repo.db.Model(&request).Where("sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?", u1, u2, u2, u1).Find(&request)
	if request.ReceiverID == 0 {
		return false
	}
	return true
}

func (repo *RequestsRepository) GetAllConnections(id uint) ([]model.Connection, []model.Connection) {
	var ids1 []model.Connection
	repo.db.Model(&ids1).Where("user_two = ?", id).Find(&ids1)
	var ids2 []model.Connection
	repo.db.Model(&ids2).Where("user_one = ?", id).Find(&ids2)

	return ids1, ids2
}

func (repo *RequestsRepository) FindMessages(id1 int64, id2 int64) []model.Message {
	var messages []model.Message
	repo.db.Model(&messages).Where("sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?", id1, id2, id2, id1).Find(&messages)

	return messages
}

func (repo *RequestsRepository) DeleteConnection(id1 int64, id2 int64) {
	repo.db.Where("user_one = ? AND user_two = ? OR user_one = ? AND user_two = ?", id1, id2, id2, id1).Delete(&model.Connection{})
	return
}

func (repo *RequestsRepository) SendNotification(id uint, text string) {
	notification := model.Notification{
		Text:       text,
		ReceiverId: id,
		Date:       time.Now(),
	}
	repo.db.Create(&notification)
}

func (repo *RequestsRepository) GetNotifications(id int64) []model.Notification {
	var notifications []model.Notification
	repo.db.Model(&notifications).Where("receiver_id = ?", id).Find(&notifications)

	return notifications
}
