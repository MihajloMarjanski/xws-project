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

func (repository *connectionRepository) CreateUserConnection( /*ctx context.Context, */ f model.Connection) (bool, error) {
	//span := tracer.StartSpanFromContextMetadata(ctx, "CreateFollowersConnection")
	//defer span.Finish()
	//ctx = tracer.ContextWithSpan(context.Background(), span)

	session := repository.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MERGE (a:User {id : $UserOne})-[:CONNECT]-(b:User {id : $UserTwo})"+
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

func (repository *connectionRepository) FindRecommendationsForUser( /*ctx context.Context, */ u model.User) ([]uint, error) {
	//span := tracer.StartSpanFromContextMetadata(ctx, "CreateFollowersConnection")
	//defer span.Finish()
	//ctx = tracer.ContextWithSpan(context.Background(), span)

	session := repository.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User {id: $UserId})-[:CONNECT]-(b:Person)-[:CONNECT]-(c:Person)"+
				"WHERE a <> c AND NOT (a)-[:CONNECT]-(c)"+
				"RETURN c.Id, count(c) as frequency"+
				"ORDER BY frequency DESC"+
				"LIMIT 5",
			map[string]interface{}{
				"UserOne": u.UserId,
			})

		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Collect(), nil
		}
		return nil, errors.New("error: no user recommendations found")
	})
	if err != nil || result == nil {
		return nil, err
	}

	return result.([]uint), nil
}
