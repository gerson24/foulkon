package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/kylelemons/godebug/pretty"
)

func TestPostgresRepo_AddGroup(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *api.Group
		// Postgres Repo Args
		groupToCreate *api.Group
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			groupToCreate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCasegroupAlreadyExist": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
			groupToCreate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"groups_pkey\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanGroupTable()

		// Insert previous data
		if test.previousGroup != nil {
			err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.UnixNano(), test.previousGroup.Urn, test.previousGroup.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
				continue
			}
		}
		// Call to repository to store group
		storedGroup, err := repoDB.AddGroup(*test.groupToCreate)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
			if diff := pretty.Compare(dbError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(storedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
			// Check database
			groupNumber, err := getGroupsCountFiltered(test.groupToCreate.ID, test.groupToCreate.Name, test.groupToCreate.Path,
				test.groupToCreate.CreateAt.UnixNano(), test.groupToCreate.Urn, test.groupToCreate.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting groups: %v", n, err)
				continue
			}
			if groupNumber != 1 {
				t.Errorf("Test %v failed. Received different group number: %v", n, groupNumber)
				continue
			}
		}
	}
}

func TestPostgresRepo_GetGroupByName(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *api.Group
		// Postgres Repo Args
		org  string
		name string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			org:  "Org",
			name: "Name",
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseGroupNotExist": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			org:  "Org",
			name: "NotExist",
			expectedError: &database.Error{
				Code:    database.GROUP_NOT_FOUND,
				Message: "Group with organization Org and name NotExist not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable()

		// Insert previous data
		if test.previousGroup != nil {
			err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.UnixNano(), test.previousGroup.Urn, test.previousGroup.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
				continue
			}
		}

		// Call to repository to get group
		receivedGroup, err := repoDB.GetGroupByName(test.org, test.name)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestPostgresRepo_GetGroupById(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *api.Group
		// Postgres Repo Args
		groupID string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			groupID: "GroupID",
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseGroupNotExist": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			groupID: "NotExist",
			expectedError: &database.Error{
				Code:    database.GROUP_NOT_FOUND,
				Message: "Group with id NotExist not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable()

		// Insert previous data
		if test.previousGroup != nil {
			err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.UnixNano(), test.previousGroup.Urn, test.previousGroup.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
				continue
			}
		}

		// Call to repository to get group
		receivedGroup, err := repoDB.GetGroupById(test.groupID)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestPostgresRepo_GetGroupsFiltered(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroups []api.Group
		// Postgres Repo Args
		org    string
		filter *api.Filter
		// Expected result
		expectedResponse []api.Group
	}{
		"OkCasePathPrefix1": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				PathPrefix: "Path",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
		},
		"OkCasePathPrefix2": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				PathPrefix: "Path123",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
			},
		},
		"OkCasePathPrefix3": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				PathPrefix: "NoPath",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{},
		},
		"OkCaseGetByOrg": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			org:    "Org1",
			filter: testFilter,
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
			},
		},
		"OkCaseGetByOrgAndPathPrefix": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			org: "Org1",
			filter: &api.Filter{
				PathPrefix: "Path123",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
			},
		},
		"OkCaseWithoutParams": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			filter: testFilter,
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org2",
				},
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable()

		// Insert previous data
		if test.previousGroups != nil {
			for _, previousGroup := range test.previousGroups {
				if err := insertGroup(previousGroup.ID, previousGroup.Name, previousGroup.Path,
					previousGroup.CreateAt.UnixNano(), previousGroup.Urn, previousGroup.Org); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous groups: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to get groups
		receivedGroups, total, err := repoDB.GetGroupsFiltered(test.org, test.filter)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check total
		if total != len(test.expectedResponse) {
			t.Errorf("Test %v failed. Received different total elements: %v", n, total)
			continue
		}
		// Check response
		if diff := pretty.Compare(receivedGroups, test.expectedResponse); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func TestPostgresRepo_UpdateGroup(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroups []api.Group
		// Postgres Repo Args
		groupToUpdate *api.Group
		newName       string
		newPath       string
		newUrn        string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn",
					CreateAt: now,
					Org:      "Org",
				},
			},
			groupToUpdate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			newName: "NewName",
			newPath: "NewPath",
			newUrn:  "NewUrn",
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "NewName",
				Path:     "NewPath",
				Urn:      "NewUrn",
				CreateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseDuplicateUrn": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn",
					CreateAt: now,
					Org:      "Org",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path2",
					Urn:      "Fail",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			groupToUpdate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			newName: "NewName",
			newPath: "NewPath",
			newUrn:  "Fail",
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"groups_urn_key\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable()

		// Insert previous data
		if test.previousGroups != nil {
			for _, previousGroup := range test.previousGroups {
				err := insertGroup(previousGroup.ID, previousGroup.Name, previousGroup.Path,
					previousGroup.CreateAt.UnixNano(), previousGroup.Urn, previousGroup.Org)
				if err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
					continue
				}
			}
		}

		// Call to repository to update group
		updatedGroup, err := repoDB.UpdateGroup(*test.groupToUpdate, test.newName, test.newPath, test.newUrn)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
			if diff := pretty.Compare(dbError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(updatedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
			// Check database
			groupNumber, err := getGroupsCountFiltered(test.expectedResponse.ID, test.expectedResponse.Name, test.expectedResponse.Path,
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.Urn, test.expectedResponse.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting groups: %v", n, err)
				continue
			}
			if groupNumber != 1 {
				t.Fatalf("Test %v failed. Received different group number: %v", n, groupNumber)
				continue
			}
		}
	}
}

func TestPostgresRepo_RemoveGroup(t *testing.T) {
	type relation struct {
		userID        string
		groupIDs      []string
		groupNotFound bool
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *api.Group
		relation      *relation
		// Postgres Repo Args
		groupToDelete string
	}{
		"OkCase": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			relation: &relation{
				userID:   "UserID",
				groupIDs: []string{"GroupID"},
			},
			groupToDelete: "GroupID",
		},
	}

	for n, test := range testcases {
		cleanGroupTable()
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.previousGroup != nil {
			if err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.Unix(), test.previousGroup.Urn, test.previousGroup.Org); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group: %v", n, err)
				continue
			}
		}
		if test.relation != nil {
			for _, id := range test.relation.groupIDs {
				if err := insertGroupUserRelation(test.relation.userID, id); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to remove group
		err := repoDB.RemoveGroup(test.groupToDelete)

		// Check database
		groupNumber, err := getGroupsCountFiltered(test.groupToDelete, "", "",
			0, "", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting groups: %v", n, err)
			continue
		}
		if groupNumber != 0 {
			t.Errorf("Test %v failed. Received different group number: %v", n, groupNumber)
			continue
		}

		relations, err := getGroupUserRelations(test.previousGroup.ID, "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
			continue
		}
	}
}

func TestPostgresRepo_AddMember(t *testing.T) {
	testcases := map[string]struct {
		// Postgres Repo Args
		userID  string
		groupID string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			userID:  "UserID",
			groupID: "GroupID",
		},
		"ErrorCaseInternalError": {
			groupID: "GroupID",
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: null value in column user_id violates not-null constraint",
			},
		},
	}

	for n, test := range testcases {
		// Clean GroupUserRelation database
		cleanGroupUserRelationTable()

		// Call to repository to store member
		err := repoDB.AddMember(test.userID, test.groupID)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}

			// Check database
			relations, err := getGroupUserRelations(test.groupID, test.userID)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
				continue
			}
			if relations != 1 {
				t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
				continue
			}
		}
	}
}

func TestPostgresRepo_RemoveMember(t *testing.T) {
	type relation struct {
		userID  string
		groupID string
	}
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		userID  string
		groupID string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			relation: &relation{
				userID:  "UserID",
				groupID: "GroupID",
			},
			userID:  "UserID",
			groupID: "GroupID",
		},
	}

	for n, test := range testcases {
		// Clean GroupUserRelation database
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupUserRelation(test.relation.userID, test.relation.groupID); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
				continue
			}
		}

		// Call to repository to remove member
		err := repoDB.RemoveMember(test.userID, test.groupID)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}

		// Check database
		relations, err := getGroupUserRelations(test.groupID, test.userID)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
			continue
		}
	}
}

func TestPostgresRepo_IsMemberOfGroup(t *testing.T) {
	type relation struct {
		userID  string
		groupID string
	}
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		group  string
		member string
		// Expected result
		isMember bool
	}{
		"OkCaseIsMember": {
			relation: &relation{
				userID:  "UserID",
				groupID: "GroupID",
			},
			group:    "GroupID",
			member:   "UserID",
			isMember: true,
		},
		"OkCaseIsNotMember": {
			group:    "GroupID",
			member:   "UserID",
			isMember: false,
		},
	}

	for n, test := range testcases {
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupUserRelation(test.relation.userID, test.relation.groupID); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
				continue
			}
		}

		isMember, err := repoDB.IsMemberOfGroup(test.member, test.group)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check response
		if diff := pretty.Compare(isMember, test.isMember); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func TestPostgresRepo_GetGroupMembers(t *testing.T) {
	type relations struct {
		users        []api.User
		groupID      string
		userNotFound bool
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relations *relations
		// Postgres Repo Args
		groupID string
		filter  *api.Filter
		// Expected result
		expectedResponse []api.User
		expectedError    *database.Error
	}{
		"OkCase": {
			relations: &relations{
				users: []api.User{
					{
						ID:         "UserID1",
						ExternalID: "ExternalID1",
						Path:       "Path",
						Urn:        "urn1",
						CreateAt:   now,
					},
					{
						ID:         "UserID2",
						ExternalID: "ExternalID2",
						Path:       "Path",
						Urn:        "urn2",
						CreateAt:   now,
					},
				},
				groupID: "GroupID",
			},
			groupID: "GroupID",
			filter:  testFilter,
			expectedResponse: []api.User{
				{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path",
					Urn:        "urn1",
					CreateAt:   now,
				},
				{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path",
					Urn:        "urn2",
					CreateAt:   now,
				},
			},
		},
		"ErrorCase": {
			relations: &relations{
				users: []api.User{
					{
						ID: "UserID1",
					},
				},
				groupID:      "GroupID",
				userNotFound: true,
			},
			groupID: "GroupID",
			filter:  testFilter,
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Code: UserNotFound, Message: User with id UserID1 not found",
			},
		},
	}

	for n, test := range testcases {
		cleanUserTable()
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.relations != nil {
			for _, user := range test.relations.users {
				if err := insertGroupUserRelation(user.ID, test.relations.groupID); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
					continue
				}
				if !test.relations.userNotFound {
					if err := insertUser(user.ID, user.ExternalID, user.Path,
						user.CreateAt.UnixNano(), user.Urn); err != nil {
						t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
						continue
					}
				}
			}

		}

		receivedUsers, total, err := repoDB.GetGroupMembers(test.groupID, test.filter)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
			if diff := pretty.Compare(dbError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check total
			if total != len(test.expectedResponse) {
				t.Errorf("Test %v failed. Received different total elements: %v", n, total)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedUsers, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestPostgresRepo_AttachPolicy(t *testing.T) {
	testcases := map[string]struct {
		// Postgres Repo Args
		policyID string
		groupID  string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			policyID: "PolicyID",
			groupID:  "GroupID",
		},
		"ErrorCaseInternalError": {
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: null value in column user_id violates not-null constraint",
			},
		},
	}

	for n, test := range testcases {
		// Clean GroupPolicyRelation database
		cleanGroupPolicyRelationTable()

		// Call to repository to attach policy
		err := repoDB.AttachPolicy(test.groupID, test.policyID)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}

			// Check database
			relations, err := getGroupPolicyRelationCount(test.policyID, test.groupID)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
				continue
			}
			if relations != 1 {
				t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
				continue
			}
		}
	}
}

func TestPostgresRepo_DetachPolicy(t *testing.T) {
	type relation struct {
		policyID string
		groupID  string
	}
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		policyID string
		groupID  string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			relation: &relation{
				policyID: "PolicyID",
				groupID:  "GroupID",
			},
			policyID: "PolicyID",
			groupID:  "GroupID",
		},
	}

	for n, test := range testcases {
		// Clean GroupPolicyRelation database
		cleanGroupPolicyRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupPolicyRelation(test.relation.groupID, test.relation.policyID); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group policy relations: %v", n, err)
				continue
			}
		}

		// Call to repository to detach policy
		err := repoDB.DetachPolicy(test.groupID, test.policyID)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}

		// Check database
		relations, err := getGroupPolicyRelationCount(test.policyID, test.groupID)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
			continue
		}
	}
}

func TestPostgresRepo_IsAttachedToGroup(t *testing.T) {
	type relation struct {
		groupID  string
		policyID string
	}
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		groupID  string
		policyID string
		// Expected result
		expectedResult bool
	}{
		"OkCase": {
			relation: &relation{
				groupID:  "GroupID",
				policyID: "PolicyID",
			},
			groupID:        "GroupID",
			policyID:       "PolicyID",
			expectedResult: true,
		},
		"OkCaseNotFound": {
			relation: &relation{
				groupID:  "GroupID",
				policyID: "PolicyID",
			},
			groupID:        "GroupID",
			policyID:       "PolicyIDXXXXXXX",
			expectedResult: false,
		},
	}

	for n, test := range testcases {
		// Clean GroupPolicyRelation database
		cleanGroupPolicyRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupPolicyRelation(test.relation.groupID, test.relation.policyID); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group policy relations: %v", n, err)
				continue
			}
		}

		// Call repository to check if policy is attached to group
		result, err := repoDB.IsAttachedToGroup(test.groupID, test.policyID)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}

		if result != test.expectedResult {
			t.Errorf("Test %v failed. Received %v, expected %v", n, result, test.expectedResult)
			continue
		}
	}
}

func TestPostgresRepo_GetAttachedPolicies(t *testing.T) {
	type relations struct {
		policies       []api.Policy
		groupID        string
		policyNotFound bool
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relations  *relations
		statements []Statement
		// Postgres Repo Args
		groupID string
		filter  *api.Filter
		// Expected result
		expectedResponse []api.Policy
		expectedError    *database.Error
	}{
		"OkCase": {
			relations: &relations{
				policies: []api.Policy{
					{
						ID:       "PolicyID1",
						Name:     "Name1",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now,
						Urn:      "Urn1",
					},
					{
						ID:       "PolicyID2",
						Name:     "Name2",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now,
						Urn:      "Urn2",
					},
				},
				groupID: "GroupID",
			},
			statements: []Statement{},
			groupID:    "GroupID",
			filter:     testFilter,
			expectedResponse: []api.Policy{
				{
					ID:         "PolicyID1",
					Name:       "Name1",
					Org:        "org1",
					Path:       "/path/",
					CreateAt:   now,
					Urn:        "Urn1",
					Statements: &[]api.Statement{},
				},
				{
					ID:         "PolicyID2",
					Name:       "Name2",
					Org:        "org1",
					Path:       "/path/",
					CreateAt:   now,
					Urn:        "Urn2",
					Statements: &[]api.Statement{},
				},
			},
		},
		"ErrorCase": {
			relations: &relations{
				policies: []api.Policy{
					{
						ID:       "PolicyID1",
						Name:     "Name1",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now,
						Urn:      "Urn1",
					},
					{
						ID:       "PolicyID2",
						Name:     "Name2",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now,
						Urn:      "Urn2",
					},
				},
				groupID:        "GroupID",
				policyNotFound: true,
			},
			statements: []Statement{},
			groupID:    "GroupID",
			filter:     testFilter,
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Code: PolicyNotFound, Message: Policy with id PolicyID1 not found",
			},
		},
	}

	for n, test := range testcases {
		cleanPolicyTable()
		cleanGroupPolicyRelationTable()

		// Insert previous data
		if test.relations != nil {
			for _, policy := range test.relations.policies {
				if err := insertGroupPolicyRelation(test.relations.groupID, policy.ID); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group policy relations: %v", n, err)
					continue
				}
				if !test.relations.policyNotFound {
					if err := insertPolicy(policy.ID, policy.Name, policy.Org, policy.Path,
						policy.CreateAt.UnixNano(), policy.Urn, test.statements); err != nil {
						t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
						continue
					}
				}
			}

		}

		receivedPolicies, total, err := repoDB.GetAttachedPolicies(test.groupID, test.filter)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
			if diff := pretty.Compare(dbError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check total
			if total != len(test.expectedResponse) {
				t.Errorf("Test %v failed. Received different total elements: %v", n, total)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedPolicies, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}
