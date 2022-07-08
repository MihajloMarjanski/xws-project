package repo

import (
	"connection-service/model"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type connectionRepository struct {
	driver neo4j.Driver
}

func NewConnectionRepository(driver neo4j.Driver) (*connectionRepository, error) {
	if driver == nil {
		panic("Connection repository not created, driver is nil")
	}

	return &connectionRepository{driver: driver}, nil
}

func (repository *connectionRepository) CreateUser( /*ctx context.Context, */ u model.User) (bool, error) {
	//span := tracer.StartSpanFromContextMetadata(ctx, "CreateFollowersConnection")
	//defer span.Finish()
	//ctx = tracer.ContextWithSpan(context.Background(), span)

	session := repository.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (n:User {id : $UserId}) RETURN n",
			map[string]interface{}{
				"UserId": u.UserId,
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

func (repository *connectionRepository) CreateUserConnection( /*ctx context.Context, */ f model.Connection) (bool, error) {
	//span := tracer.StartSpanFromContextMetadata(ctx, "CreateFollowersConnection")
	//defer span.Finish()
	//ctx = tracer.ContextWithSpan(context.Background(), span)

	session := repository.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User {id : $UserOne}), (b:User {id : $UserTwo}) MERGE"+
				"(a)-[:CONNECT]-(b)"+
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
