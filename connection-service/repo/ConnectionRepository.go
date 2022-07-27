package repo

import (
	"connection-service/model"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ConnectionRepository struct {
	driver neo4j.Driver
}

func NewConnectionRepository() *ConnectionRepository {
	driver, err := neo4j.NewDriver("neo4j://neo4j:7687", neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		panic("Connection repository not created, driver is nil")
	}
	return &ConnectionRepository{driver: driver}
}

func (repository *ConnectionRepository) CreateUser( /*ctx context.Context, */ id uint) (bool, error) {
	//span := tracer.StartSpanFromContextMetadata(ctx, "CreateFollowersConnection")
	//defer span.Finish()
	//ctx = tracer.ContextWithSpan(context.Background(), span)

	session := repository.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		t := fmt.Sprintf("u create: %d", id)
		fmt.Println(t)

		result, err := transaction.Run(
			"CREATE (n:User {id : $UserId}) RETURN n",
			map[string]interface{}{
				"UserId": id,
			})

		if err != nil {
			return nil, err
		}
		if result.Next() {
			return true, nil
		}
		return false, errors.New("error: can not create user ")
	})
	if err != nil || result == nil {
		return false, err
	}
	return true, nil
}

func checkIfUserExist(id uint, transaction neo4j.Transaction) bool {
	t := fmt.Sprintf("dobio za chechUserExist: %d", id)
	fmt.Println(t)

	result, _ := transaction.Run(
		"MATCH (a:User { id: $UserId }) RETURN a",
		map[string]interface{}{"UserId": id})

	if result != nil && result.Next() {
		//t := fmt.Sprintf(result)
		//fmt.Println(t)
		return true
	}
	return false
}

func (repository *ConnectionRepository) CreateUserConnection( /*ctx context.Context, */ f model.Connection) (bool, error) {
	//span := tracer.StartSpanFromContextMetadata(ctx, "CreateFollowersConnection")
	//defer span.Finish()
	//ctx = tracer.ContextWithSpan(context.Background(), span)

	session := repository.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		if !checkIfUserExist(f.UserOne, transaction) {
			repository.CreateUser(f.UserOne)
		}

		if !checkIfUserExist(f.UserTwo, transaction) {
			repository.CreateUser(f.UserTwo)
		}

		result, err := transaction.Run(
			"MATCH (a:User {id : $UserOne}), (b:User {id : $UserTwo}) MERGE (a)-[:CONNECT]-(b)"+
				"RETURN a, b",
			map[string]interface{}{
				"UserOne": f.UserOne,
				"UserTwo": f.UserTwo,
			})

		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Record().Values[0], nil
		}
		return false, errors.New("error: can not create connection")
	})
	if err != nil || result == nil {
		return false, err
	}
	return true, nil
}

func (repository *ConnectionRepository) FindRecommendationsForUser( /*ctx context.Context, */ u model.User) ([]uint, error) {
	//span := tracer.StartSpanFromContextMetadata(ctx, "CreateFollowersConnection")
	//defer span.Finish()
	//ctx = tracer.ContextWithSpan(context.Background(), span)
	var list []uint

	session := repository.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User {id: $UserId})-[:CONNECT]-(b:User)-[:CONNECT]-(c:User)"+
				"WHERE a <> c AND NOT (a)-[:CONNECT]-(c)"+
				"RETURN c.id, count(c) as frequency"+
				"ORDER BY frequency DESC"+
				"LIMIT 5",
			map[string]interface{}{
				"UserOne": u.UserId,
			})

		if err != nil {
			return nil, err
		}
		if result.Next() {
			ids, _ := result.Collect()
			for _, userId := range ids {
				fmt.Println(userId)
				list = append(list, userId.Values[0].(uint))
			}
			fmt.Println(list)
			return list, nil
		}
		return nil, errors.New("error: no user recommendations found")
	})
	if err != nil || result == nil {
		return nil, err
	}

	return list, nil
}
