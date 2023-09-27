package models_tests

import (
	"bullshape/db"
	"bullshape/models"
	"bullshape/utils/logger"
	"net/http"
	"testing"
)

func TestGetCompany(t *testing.T) {

	type test struct {
		name           string
		companiesToAdd []db.Company
		expectedToFind bool
	}
	tt := []test{
		{
			name: "Get Company - Successful",
			companiesToAdd: []db.Company{
				{
					Name:        "Name1",
					UUID:        "UUID-1",
					Description: "Description",
					NumEmployes: 1,
					Registered:  true,
					Type:        models.COOPERATIVE,
				},
			},
			expectedToFind: true,
		},
		{
			name: "Get Company - No result",
			companiesToAdd: []db.Company{
				{
					Name:        "Name1",
					UUID:        "UUID-1",
					Description: "Description",
					NumEmployes: 1,
					Registered:  true,
					Type:        models.COOPERATIVE,
				},
			},
			expectedToFind: false,
		},
	}
	for _, test := range tt {

		t.Run(test.name, func(t *testing.T) {
			log := logger.GetLogger()
			db := db.NewDatabaseConnection()
			m := models.NewCtrlServices(log.Slogger, db)

			//prepare
			for i := range test.companiesToAdd {
				test.companiesToAdd[i].Create(m.DB)
				//cleanup
				defer test.companiesToAdd[i].Delete(m.DB)
			}

			//test
			var err error
			var company *models.Company
			if test.expectedToFind {
				company, _, err = m.GetCompany(test.companiesToAdd[0].ID)
			} else {
				max := len(test.companiesToAdd)
				company, _, err = m.GetCompany(test.companiesToAdd[max-1].ID + 1)
			}

			//assert
			if test.expectedToFind && err != nil {
				t.Fatal("Expected no error. Got error: ", err)
			}

			if !test.expectedToFind && err == nil {
				t.Fatal("Expected  error. Got no error: ")
			}
			if test.expectedToFind {
				if test.companiesToAdd[0].Name != company.Name {
					t.Fatal("Wrong company name. Expected : ", test.companiesToAdd[0].Name,
						"Got: ", company.Name)
				}
				if test.companiesToAdd[0].UUID != company.UUID {
					t.Fatal("Wrong company uuid. Expected : ", test.companiesToAdd[0].UUID,
						"Got: ", company.UUID)
				}
				if test.companiesToAdd[0].Description != company.Description {
					t.Fatal("Wrong company description. Expected : ", test.companiesToAdd[0].Description,
						"Got: ", company.Description)
				}
				if test.companiesToAdd[0].NumEmployes != company.NumEmployes {
					t.Fatal("Wrong company number of employes. Expected : ", test.companiesToAdd[0].NumEmployes,
						"Got: ", company.NumEmployes)
				}
				if test.companiesToAdd[0].Registered != company.Registered {
					t.Fatal("Wrong company register status. Expected : ", test.companiesToAdd[0].Registered,
						"Got: ", company.Registered)
				}
				if test.companiesToAdd[0].Type != company.Type {
					t.Fatal("Wrong company name. Expected : ", test.companiesToAdd[0].Type,
						"Got: ", company.Type)
				}
			}

		})

	}

}

func TestCreateCompany(t *testing.T) {

	type test struct {
		name           string
		companiesToAdd models.NewCompany
		expectedError  bool
		expectedCode   int
	}
	name := "name1"
	uuid := "uuid"
	numEmpl := uint(1)
	registered := true
	tp, wrongTp := db.COOPERATIVE, "mytype"
	tt := []test{
		{
			name: "Create Company - Successful",
			companiesToAdd: models.NewCompany{

				Name:        &name,
				UUID:        &uuid,
				Description: "Description",
				NumEmployes: &numEmpl,
				Registered:  &registered,
				Type:        &tp,
			},
			expectedError: false,
			expectedCode:  http.StatusOK,
		},

		{
			name: "Create Company - Fail - Empty Name",
			companiesToAdd: models.NewCompany{
				Name:        nil,
				UUID:        &uuid,
				Description: "Description",
				NumEmployes: &numEmpl,
				Registered:  &registered,
				Type:        &tp,
			},
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name: "Create Company - Fail - Empty NumEmployes",
			companiesToAdd: models.NewCompany{

				Name:        &name,
				UUID:        &uuid,
				Description: "Description",
				NumEmployes: nil,
				Registered:  &registered,
				Type:        &tp,
			},
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name: "Create Company - Fail - Empty type",
			companiesToAdd: models.NewCompany{

				Name:        &name,
				UUID:        &uuid,
				Description: "Description",
				NumEmployes: &numEmpl,
				Registered:  &registered,
				Type:        nil,
			},
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name: "Create Company - Fail - Invalid type",
			companiesToAdd: models.NewCompany{

				Name:        &name,
				UUID:        &uuid,
				Description: "Description",
				NumEmployes: &numEmpl,
				Registered:  &registered,
				Type:        &wrongTp,
			},
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
	}
	for _, test := range tt {

		t.Run(test.name, func(t *testing.T) {
			log := logger.GetLogger()
			dbc := db.NewDatabaseConnection()
			m := models.NewCtrlServices(log.Slogger, dbc)

			//test
			var err error
			var company *models.Company
			var status int

			company, status, err = m.CreateCompany(test.companiesToAdd)
			var dbCompany *db.Company
			var errDB error

			if company != nil {
				dbCompany, errDB = db.GetCompanyByID(m.DB, company.ID)
				if errDB == nil {
					defer dbCompany.Delete(m.DB)
				}
			}

			//assert
			if test.expectedError && err == nil {
				t.Fatal("Expected  error. Got no error ")
			}
			if !test.expectedError && err != nil {
				t.Fatal("Expected  no error. Got error: ", err)
			}
			if test.expectedCode != status {
				t.Fatal("Expected status: ", test.expectedCode, ". Got: ", status)
			}

			if !test.expectedError {
				if errDB != nil || dbCompany == nil {
					t.Fatal("Expected a new company DB enty Got nothing. Error : ", errDB)
				}
				defer dbCompany.Delete(m.DB)
				if dbCompany.Name != company.Name {
					t.Fatal("Wrong company name. Expected : ", dbCompany.Name,
						"Got: ", company.Name)
				}
				if dbCompany.UUID != company.UUID {
					t.Fatal("Wrong company uuid. Expected : ", dbCompany.UUID,
						"Got: ", company.UUID)
				}
				if dbCompany.Description != company.Description {
					t.Fatal("Wrong company description. Expected : ", dbCompany.Description,
						"Got: ", company.Description)
				}
				if dbCompany.NumEmployes != company.NumEmployes {
					t.Fatal("Wrong company number of employes. Expected : ", dbCompany.NumEmployes,
						"Got: ", company.NumEmployes)
				}
				if dbCompany.Registered != company.Registered {
					t.Fatal("Wrong company register status. Expected : ", dbCompany.Registered,
						"Got: ", company.Registered)
				}
				if dbCompany.Type != company.Type {
					t.Fatal("Wrong company name. Expected : ", dbCompany.Type,
						"Got: ", company.Type)
				}
			}

		})

	}

}

func TestCreateTwoCompanies(t *testing.T) {

	type test struct {
		name           string
		companiesToAdd []models.NewCompany
		expectedError  bool
		expectedCode   int
	}
	name1, name2 := "name1", "name2"
	uuid1, uuid2 := "uuid", "uuid2"
	numEmpl := uint(1)
	registered := true
	tp := db.COOPERATIVE
	tt := []test{
		{
			name: "Create two Companies - Successful",
			companiesToAdd: []models.NewCompany{
				{
					Name:        &name1,
					UUID:        &uuid1,
					Description: "Description",
					NumEmployes: &numEmpl,
					Registered:  &registered,
					Type:        &tp,
				},
				{
					Name:        &name2,
					UUID:        &uuid2,
					Description: "Description",
					NumEmployes: &numEmpl,
					Registered:  &registered,
					Type:        &tp,
				},
			},
			expectedError: false,
			expectedCode:  http.StatusOK,
		},
		{
			name: "Create two Companies - Fail - Companies Same Name",
			companiesToAdd: []models.NewCompany{
				{
					Name:        &name1,
					UUID:        &uuid1,
					Description: "Description",
					NumEmployes: &numEmpl,
					Registered:  &registered,
					Type:        &tp,
				},
				{
					Name:        &name1,
					UUID:        &uuid2,
					Description: "Description",
					NumEmployes: &numEmpl,
					Registered:  &registered,
					Type:        &tp,
				},
			},
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
	}
	for _, test := range tt {

		t.Run(test.name, func(t *testing.T) {
			log := logger.GetLogger()
			dbc := db.NewDatabaseConnection()
			m := models.NewCtrlServices(log.Slogger, dbc)

			//test
			var err error
			var company *models.Company
			var status int
			for _, companyToAdd := range test.companiesToAdd {
				company, status, err = m.CreateCompany(companyToAdd)
				if err != nil {
					break
				}
			}
			//cleanup
			dbCompanies, _ := db.GetCompanies(m.DB)
			for i := range dbCompanies {
				defer dbCompanies[i].Delete(m.DB)
			}

			//assert
			if test.expectedError && err == nil {
				t.Fatal("Expected  error. Got no error ")
			}
			if !test.expectedError && err != nil {
				t.Fatal("Expected  no error. Got error: ", err)
			}
			if test.expectedCode != status {
				t.Fatal("Expected status: ", test.expectedCode, ". Got: ", status)
			}

			if !test.expectedError {
				dbCompany, err := db.GetCompanyByID(m.DB, company.ID)
				if err != nil || dbCompany == nil {
					t.Fatal("Expected a new company DB enty Got nothing. Error : ", err)
				}
				if dbCompany.Name != company.Name {
					t.Fatal("Wrong company name. Expected : ", dbCompany.Name,
						"Got: ", company.Name)
				}
				if dbCompany.UUID != company.UUID {
					t.Fatal("Wrong company uuid. Expected : ", dbCompany.UUID,
						"Got: ", company.UUID)
				}
				if dbCompany.Description != company.Description {
					t.Fatal("Wrong company description. Expected : ", dbCompany.Description,
						"Got: ", company.Description)
				}
				if dbCompany.NumEmployes != company.NumEmployes {
					t.Fatal("Wrong company number of employes. Expected : ", dbCompany.NumEmployes,
						"Got: ", company.NumEmployes)
				}
				if dbCompany.Registered != company.Registered {
					t.Fatal("Wrong company register status. Expected : ", dbCompany.Registered,
						"Got: ", company.Registered)
				}
				if dbCompany.Type != company.Type {
					t.Fatal("Wrong company name. Expected : ", dbCompany.Type,
						"Got: ", company.Type)
				}
			}

		})

	}

}
func TestDeleteCompany(t *testing.T) {

	type test struct {
		name                string
		companiesToAdd      db.Company
		expectedToBeDeleted bool
	}
	tt := []test{
		{
			name: "Delete Company - Successful",
			companiesToAdd: db.Company{

				Name:        "Name1",
				UUID:        "UUID-1",
				Description: "Description",
				NumEmployes: 1,
				Registered:  true,
				Type:        models.COOPERATIVE,
			},
			expectedToBeDeleted: true,
		},
		{
			name: "Delete Company - wrong ID. No action",
			companiesToAdd: db.Company{

				Name:        "Name1",
				UUID:        "UUID-1",
				Description: "Description",
				NumEmployes: 1,
				Registered:  true,
				Type:        models.COOPERATIVE,
			},
			expectedToBeDeleted: false,
		},
	}
	for _, test := range tt {

		t.Run(test.name, func(t *testing.T) {
			log := logger.GetLogger()
			dbc := db.NewDatabaseConnection()
			m := models.NewCtrlServices(log.Slogger, dbc)

			//prepare
			test.companiesToAdd.Create(m.DB)
			//cleanup
			defer test.companiesToAdd.Delete(m.DB)

			//test
			var err error
			if test.expectedToBeDeleted {
				_, err = m.DeleteCompany(test.companiesToAdd.ID)
			} else {
				_, err = m.DeleteCompany(test.companiesToAdd.ID + 1)
			}

			//assert
			if test.expectedToBeDeleted && err != nil {
				t.Fatal("Expected no error. Got error: ", err)
			}

			if !test.expectedToBeDeleted && err == nil {
				t.Fatal("Expected  error. Got no error: ")
			}
			companyDB, err := db.GetCompanyByID(m.DB, test.companiesToAdd.ID)
			if test.expectedToBeDeleted && err == nil {
				t.Fatal("Expected no company DB entry. Got DB entry with ID: ", companyDB.ID)
			}

			if !test.expectedToBeDeleted && err != nil {
				t.Fatal("Expected a company DB entry. Got no DB entry")
			}

		})

	}

}

func TestUpdateCompany(t *testing.T) {

	type test struct {
		name           string
		companiesToAdd db.Company
		companyopts    models.EditCompanyOpts
		expectedError  bool
		expectedStatus int
	}
	newDescription := "New Description"
	newNumberOfEmployes := uint(1)
	newRegisterd := false
	newType, newWrongType := db.COOPERATIVE, "mytype"
	tt := []test{
		{
			name: "Update Company - All fields -Successful",
			companiesToAdd: db.Company{

				Name:        "Name1",
				UUID:        "UUID-1",
				Description: "Description",
				NumEmployes: 100,
				Registered:  true,
				Type:        models.COOPERATIVE,
			},
			companyopts: models.EditCompanyOpts{
				Description: &newDescription,
				NumEmployes: &newNumberOfEmployes,
				Registered:  &newRegisterd,
				Type:        &newType,
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Update Company - Only Desctiption - Successful",
			companiesToAdd: db.Company{

				Name:        "Name1",
				UUID:        "UUID-1",
				Description: "Description",
				NumEmployes: 100,
				Registered:  true,
				Type:        models.COOPERATIVE,
			},
			companyopts: models.EditCompanyOpts{
				Description: &newDescription,
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Update Company - Only Num of Employes - Successful",
			companiesToAdd: db.Company{

				Name:        "Name1",
				UUID:        "UUID-1",
				Description: "Description",
				NumEmployes: 100,
				Registered:  true,
				Type:        models.COOPERATIVE,
			},
			companyopts: models.EditCompanyOpts{
				NumEmployes: &newNumberOfEmployes,
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Update Company - No update. Wrong Type -Successful",
			companiesToAdd: db.Company{

				Name:        "Name1",
				UUID:        "UUID-1",
				Description: "Description",
				NumEmployes: 100,
				Registered:  true,
				Type:        models.COOPERATIVE,
			},
			companyopts: models.EditCompanyOpts{
				Type: &newWrongType,
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
	}
	for _, test := range tt {

		t.Run(test.name, func(t *testing.T) {
			log := logger.GetLogger()
			dbc := db.NewDatabaseConnection()
			m := models.NewCtrlServices(log.Slogger, dbc)

			//prepare
			test.companiesToAdd.Create(m.DB)
			//cleanup
			defer test.companiesToAdd.Delete(m.DB)

			//test
			var status int
			var err error
			if !test.expectedError {
				_, status, err = m.UpdateCompany(test.companiesToAdd.ID, test.companyopts)
			} else {
				_, status, err = m.UpdateCompany(test.companiesToAdd.ID+1, test.companyopts)
			}

			//assert
			if test.expectedError && err == nil {
				t.Fatal("Expected error. Got no error: ")
			}

			if !test.expectedError && err != nil {
				t.Fatal("Expected no error. Got error: ", err)
			}
			if test.expectedStatus != status {
				t.Fatal("Expected status: ", test.expectedStatus, " Got: ", status)
			}
			companyDB, err := db.GetCompanyByID(m.DB, test.companiesToAdd.ID)
			if !test.expectedError {
				if test.companyopts.Description != nil {
					if companyDB.Description != *test.companyopts.Description {
						t.Fatal("Expected company Description to be updated")
					}
				}
				if test.companyopts.NumEmployes != nil {
					if companyDB.NumEmployes != *test.companyopts.NumEmployes {
						t.Fatal("Expected company number of employes to be updated")
					}
				}
				if test.companyopts.Registered != nil {
					if companyDB.Registered != *test.companyopts.Registered {
						t.Fatal("Expected company registered value to be updated")
					}
				}
				if test.companyopts.Type != nil {
					if models.IsValidType(*test.companyopts.Type) {

						if companyDB.Type != *test.companyopts.Type {
							t.Fatal("Expected company type to be updated")
						}
					} else {
						if companyDB.Type == *test.companyopts.Type {
							t.Fatal("Expected company type not to be updated")
						}
					}
				}
			}
		})

	}

}
