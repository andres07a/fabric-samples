package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract proporciona funciones para gestionar un usuario
type SmartContract struct {
	contractapi.Contract
}

// User describe information básica de un usuario
type User struct {
	Key                    string `json:"key"`
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	VaccineSchemaCompleted string `json:"vaccineSchemaCompleted,omitempty" metadata:",optional"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"key"`
	Record *User
}

// InitLedger adds a base set of users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []User{
		User{ID: "1", Name: "Alan"},
		User{ID: "2", Name: "Pedro"},
		User{ID: "3", Name: "María"},
	}

	for _, user := range users {
		userAsBytes, _ := json.Marshal(user)
		err := ctx.GetStub().PutState(user.ID, userAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Create añade un nuevo usuario al 'world state' con los detalles dados
func (s *SmartContract) Create(ctx contractapi.TransactionContextInterface, id string, name string) error {
	user := User{
		ID:   id,
		Name: name,
	}

	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(user.ID, userAsBytes)
}

// FindOne returns the user stored in the world state with given key
func (s *SmartContract) FindOne(ctx contractapi.TransactionContextInterface, key string) (*User, error) {
	userAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", key)
	}

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}

// func (s *SmartContract) FindOneWithSchema(ctx contractapi.TransactionContextInterface, key string) (*peer.Response, error) {
// 	userAsBytes, err := ctx.GetStub().GetState(key)

// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
// 	}

// 	if userAsBytes == nil {
// 		return nil, fmt.Errorf("%s does not exist", key)
// 	}

// 	user := new(User)
// 	_ = json.Unmarshal(userAsBytes, user)
// 	fmt.Println("InvokeChaincode")

// 	// TODO: DA PROBLEMAS, DEPURAR
// 	response := ctx.GetStub().InvokeChaincode("vaccine",
// 		[][]byte{
// 			[]byte("QueryVaccineById"),
// 			[]byte(user.ID),
// 		},
// 		"salud-estadio")
// 	fmt.Println("response")
// 	fmt.Println(response)
// 	if response.Status != http.StatusOK {
// 		return nil, errors.New(response.Message)
// 	}
// 	// return user, nil
// 	return &response, nil
// }

// FindAll returns all user found in world state
func (s *SmartContract) FindAll(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		queryResult := QueryResult{Key: queryResponse.Key, Record: user}
		results = append(results, queryResult)
	}

	return results, nil
}

// Update updates the name of user with given id in world state
func (s *SmartContract) Update(ctx contractapi.TransactionContextInterface, userId string, userName string) error {
	user, err := s.FindOne(ctx, userId)
	if err != nil {
		return err
	}

	user.Name = userName

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(userId, userAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create user chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting user chaincode: %s", err.Error())
	}
}
