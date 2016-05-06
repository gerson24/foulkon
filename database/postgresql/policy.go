package postgresql

import (
	"fmt"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

func (p PostgresRepo) GetPolicyByName(org string, name string) (*api.Policy, error) {
	policy := &Policy{}
	query := p.Dbmap.Where("org like ? AND name like ?", org, name).First(policy)

	// Check if policy exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.POLICY_NOT_FOUND,
			Message: fmt.Sprintf("Policy with organization %v and name %v not found", org, name),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Retrieve associated statements
	statements := []Statement{}
	query = p.Dbmap.Where("policy_id like ?", policy.ID).Find(&statements)
	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create API policy
	policyApi := policyDBToPolicyAPI(policy)
	policyApi.Statements = statementsDBToStatetmentsAPI(statements)

	// Return policy
	return policyApi, nil
}

func (p PostgresRepo) AddPolicy(policy api.Policy) (*api.Policy, error) {
	// Create policy model
	policyDB := &Policy{
		ID:       policy.ID,
		Name:     policy.Name,
		Path:     policy.Path,
		CreateAt: time.Now().UTC().UnixNano(),
		Urn:      policy.Urn,
		Org:      policy.Org,
	}

	transaction := p.Dbmap.Begin()

	// Create policy
	transaction.Create(policyDB)

	// Create statements
	for _, statementApi := range *policy.Statements {
		// Create statement model
		statementDB := &Statement{
			ID:        uuid.NewV4().String(),
			PolicyID:  policy.ID,
			Effect:    statementApi.Effect,
			Action:    stringArrayToSplitedString(statementApi.Action),
			Resources: stringArrayToSplitedString(statementApi.Resources),
		}
		transaction.Create(statementDB)
	}

	// Error handling
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	} else {
		transaction.Commit()
	}

	// Create API policy
	policyApi := policyDBToPolicyAPI(policyDB)
	policyApi.Statements = policy.Statements

	return policyApi, nil
}

// Private helper methods

// Transform a policy retrieved from db into a policy for API
func policyDBToPolicyAPI(policydb *Policy) *api.Policy {
	return &api.Policy{
		ID:       policydb.ID,
		Name:     policydb.Name,
		Path:     policydb.Path,
		CreateAt: time.Unix(0, policydb.CreateAt).UTC(),
		Urn:      policydb.Urn,
		Org:      policydb.Org,
	}
}

// Transform a list of statements from db into API statements
func statementsDBToStatetmentsAPI(statements []Statement) *[]api.Statement {
	statementsApi := make([]api.Statement, len(statements), cap(statements))
	for i, s := range statements {
		statementsApi[i] = api.Statement{
			Action:    strings.Split(s.Action, ";"),
			Effect:    s.Effect,
			Resources: strings.Split(s.Resources, ";"),
		}
	}

	return &statementsApi
}

// Transform an array of strings into a separated string values
func stringArrayToSplitedString(array []string) string {
	stringVal := ""
	for _, s := range array {
		if len(stringVal) == 0 {
			stringVal = s
		} else {
			stringVal = stringVal + ";" + s
		}
	}

	return stringVal
}
