package service

//
//import (
//	"GoMJTrainingCamp/dbs/dbConnection"
//	"GoMJTrainingCamp/dbs/models"
//	"github.com/DATA-DOG/go-sqlmock"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/assert"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//// Setup TestDB dengan sqlmock
//func SetupTestDB() (*gorm.DB, sqlmock.Sqlmock) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		panic("Failed to create sqlmock")
//	}
//
//	dialector := mysql.New(mysql.Config{
//		Conn:                      db,
//		SkipInitializeWithVersion: true,
//	})
//
//	testDB, err := gorm.Open(dialector, &gorm.Config{})
//	if err != nil {
//		panic("Failed to open gorm DB with sqlmock")
//	}
//
//	// Assign test DB ke dbConnection
//	dbConnection.DB = testDB
//	return testDB, mock
//}
//
//// Unit test untuk CreateTrainingClass
//func TestCreateTrainingClass(t *testing.T) {
//	db, mock := SetupTestDB()
//	defer db.DB()
//
//	service := NewClassService()
//
//	// Mock input data
//	class := &models.TrainingClass{
//		ClassName:        "Yoga",
//		ClassRequirement: "Bring your own mat",
//		ClassDateTime:    time.Now(),
//		ClassCapacity:    10,
//	}
//
//	// Mock DB query
//	mock.ExpectBegin()
//	mock.ExpectExec("INSERT INTO `training_classes`").
//		WithArgs(class.ClassName, class.ClassRequirement, class.ClassDateTime, class.ClassCapacity).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	// Execute method
//	err := service.CreateTrainingClass(class)
//	assert.NoError(t, err)
//}
//
//// Unit test untuk GetClasses dengan date
//func TestGetClassesByDate(t *testing.T) {
//	db, mock := SetupTestDB()
//	defer db.DB()
//
//	service := NewClassService()
//	date := "2023-11-14"
//
//	// Mock input data
//	class := models.TrainingClass{
//		ID:            1,
//		ClassName:     "Pilates",
//		ClassDateTime: time.Now(),
//		ClassCapacity: 20,
//	}
//
//	// Mock DB query
//	rows := sqlmock.NewRows([]string{"id", "class_name", "class_date_time", "class_capacity"}).
//		AddRow(class.ID, class.ClassName, class.ClassDateTime, class.ClassCapacity)
//
//	mock.ExpectQuery("SELECT \\* FROM `training_classes` WHERE DATE\\(class_date_time\\) = ?").
//		WithArgs(date).
//		WillReturnRows(rows)
//
//	// Execute method
//	classes, err := service.GetClasses("", date)
//	assert.NoError(t, err)
//	assert.Equal(t, 1, len(classes))
//	assert.Equal(t, "Pilates", classes[0].ClassName)
//}
//
//// Unit test untuk GetClasses tanpa filter
//func TestGetAllClasses(t *testing.T) {
//	db, mock := SetupTestDB()
//	defer db.DB()
//
//	service := NewClassService()
//
//	// Mock input data
//	class1 := models.TrainingClass{
//		ID:            1,
//		ClassName:     "Zumba",
//		ClassDateTime: time.Now(),
//		ClassCapacity: 15,
//	}
//
//	class2 := models.TrainingClass{
//		ID:            2,
//		ClassName:     "HIIT",
//		ClassDateTime: time.Now(),
//		ClassCapacity: 25,
//	}
//
//	// Mock DB query
//	rows := sqlmock.NewRows([]string{"id", "class_name", "class_date_time", "class_capacity"}).
//		AddRow(class1.ID, class1.ClassName, class1.ClassDateTime, class1.ClassCapacity).
//		AddRow(class2.ID, class2.ClassName, class2.ClassDateTime, class2.ClassCapacity)
//
//	mock.ExpectQuery("SELECT \\* FROM `training_classes`").
//		WillReturnRows(rows)
//
//	// Execute method
//	classes, err := service.GetClasses("", "")
//	assert.NoError(t, err)
//	assert.Equal(t, 2, len(classes))
//	assert.Equal(t, "Zumba", classes[0].ClassName)
//	assert.Equal(t, "HIIT", classes[1].ClassName)
//}
