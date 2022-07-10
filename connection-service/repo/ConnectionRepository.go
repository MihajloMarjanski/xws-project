package repo

import (
	"connection-service/model"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ConnectionRepository struct {
	driver neo4j.Driver
}

func NewConnectionRepository() *ConnectionRepository {
	driver, err := neo4j.NewDriver("neo4j://neo4j:7474", neo4j.BasicAuth("neo4j", "neo4j", ""))
	if err == nil {
		panic("Connection repository not created, driver is nil")
	}
	return &ConnectionRepository{driver: driver}
}

func (repository *ConnectionRepository) CreateUserConnection( /*ctx context.Context, */ f model.Connection) (bool, error) {
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

func (repository *ConnectionRepository) FindRecommendationsForUser( /*ctx context.Context, */ u model.User) ([]uint64, error) {
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
			return result.Collect()
		}
		return nil, errors.New("error: no user recommendations found")
	})
	if err != nil || result == nil {
		return nil, err
	}

	return result.([]uint64), nil
}
