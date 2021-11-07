package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract proporciona funciones para gestionar una vacuna
type SmartContract struct {
	contractapi.Contract
}

// Vaccine describe information básica de una vacuna
type Vaccine struct {
	Key    string `json:"key"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Dose   string `json:"dose"`
	Scheme string `json:"scheme"`
}

// QueryResult estructura utilizada para manejar el resultado de la consulta
type QueryResult struct {
	Key    string `json:"key"`
	Record *Vaccine
}

// InitLedger agrega un conjunto básico de vacunas al libro mayor
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	vaccines := []Vaccine{
		// usuario 1 con dos vacunas, esquema completo
		Vaccine{ID: "1", Name: "astrazeneca", Dose: "1", Scheme: "sch1", Key: "1-1-sch1"},
		Vaccine{ID: "1", Name: "astrazeneca", Dose: "2", Scheme: "sch1", Key: "1-2-sch1"},

		// usuario 2 con una vacuna, esquema incompleto
		Vaccine{ID: "2", Name: "astrazeneca", Dose: "1", Scheme: "sch1", Key: "2-1-sch1"},
	}

	for _, vaccine := range vaccines {
		vaccineAsBytes, _ := json.Marshal(vaccine)
		err := ctx.GetStub().PutState(vaccine.Key, vaccineAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Create añade una nueva vacuna al 'world state' con los detalles dados
func (s *SmartContract) Create(ctx contractapi.TransactionContextInterface, key, id, name, dose, scheme string) error {
	vaccine := Vaccine{
		Key:    key,
		ID:     id,
		Name:   name,
		Dose:   dose,
		Scheme: scheme,
	}

	vaccineAsBytes, _ := json.Marshal(vaccine)

	return ctx.GetStub().PutState(vaccine.Key, vaccineAsBytes)
}

// FindOne retorna la vacuna almacenada en el 'world state' por la key
func (s *SmartContract) FindOne(ctx contractapi.TransactionContextInterface, key string) (*Vaccine, error) {
	vaccineAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if vaccineAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", key)
	}

	vaccine := new(Vaccine)
	_ = json.Unmarshal(vaccineAsBytes, vaccine)

	return vaccine, nil
}

// FindAll retorna todas las vacunas almacenadas en el 'world state'
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

		vaccine := new(Vaccine)
		_ = json.Unmarshal(queryResponse.Value, vaccine)

		queryResult := QueryResult{Key: queryResponse.Key, Record: vaccine}
		results = append(results, queryResult)
	}

	return results, nil
}

// Update modifica los campos de una vacuna por la key proporcionada en el 'world state'
func (s *SmartContract) Update(ctx contractapi.TransactionContextInterface, vaccineKey string, vaccineName string, vaccineDose string, vaccineScheme string) error {
	vaccine, err := s.FindOne(ctx, vaccineKey)
	if err != nil {
		return err
	}

	if vaccineName != "" {
		vaccine.Name = vaccineName
	}
	if vaccineDose != "" {
		vaccine.Dose = vaccineDose
	}
	if vaccineScheme != "" {
		vaccine.Scheme = vaccineScheme
	}

	vaccineAsBytes, err := json.Marshal(vaccine)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(vaccineKey, vaccineAsBytes)
}

// FindOneById retorna las vacunas que coinciden por ID
func (s *SmartContract) FindOneById(ctx contractapi.TransactionContextInterface, id string) ([]*Vaccine, error) {
	queryString := fmt.Sprintf(`{"selector":{"id":"%v"}}`, id)

	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// getQueryResultForQueryString executes the passed in query string.
func (s *SmartContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Vaccine, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("vaccine", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []*Vaccine{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var vaccine *Vaccine

		err = json.Unmarshal(response.Value, &vaccine)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		results = append(results, vaccine)
	}
	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create vaccine chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting vaccine chaincode: %s", err.Error())
	}
}
