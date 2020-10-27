package repositories

import (
	"internship_project/models"
	"internship_project/utils"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddEar(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(EarRepo.DB)
		defer SetupTables(EarRepo.DB)
		err := EarRepo.AddEar(&testEar)
		assert.Error(err, "Error was not thrown while inserting in non-existing table")
	})

	t.Run("successful query", func(t *testing.T) {
		oldEars, _ := EarRepo.GetAllEars()
		err := EarRepo.AddEar(&testEar)
		newEars, _ := EarRepo.GetAllEars()

		assert.NoError(err)
		assert.Equal(1, len(newEars)-len(oldEars), "Ear was not added.")

	})

}

func TestGetAllEars(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful query", func(t *testing.T) {
		allEars, err := EarRepo.GetAllEars()

		assert.NoError(err, "Error was thrown while reading ears")
		assert.NotNil(allEars, "Ears returned are nil.")
		assert.IsType(allEars, []models.ExternalRights{}, "Returned result is not of type Ear")

	})
}

func TestGetEar(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(EarRepo.DB)
		defer SetupTables(EarRepo.DB)
		_, err := EarRepo.GetEar(uuid.NewV4().String())
		assert.Error(err, "Error was not thrown while getting from non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		_, err := EarRepo.GetEar(uuid)
		assert.Error(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		_, err := EarRepo.GetEar(uuid)
		assert.Error(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		EarRepo.AddEar(&testEar)
		ear, err := EarRepo.GetEar(testEar.ID)
		assert.NotNil(ear, "Result is nil")
		assert.NoError(err, "There was error while getting ear")
		assert.Equal(testEar.ID, ear.ID, "Returned ear ID and test ID do not match.")
	})
}

func TestUpdateEar(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(EarRepo.DB)
		defer SetupTables(EarRepo.DB)
		err := EarRepo.UpdateEar(testEar)
		assert.Error(err, "Error was not thrown while updating in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		testEar.ID = uuid
		err := EarRepo.UpdateEar(testEar)
		assert.NotNil(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		testEar.ID = uuid
		err := EarRepo.UpdateEar(testEar)
		assert.NotNil(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		EarRepo.AddEar(&testEar)
		testEar.Delete = true
		err := EarRepo.UpdateEar(testEar)
		assert.NoError(err, "Ear was not updated.")
	})

}

func TestDeleteEar(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(EarRepo.DB)
		defer SetupTables(EarRepo.DB)
		err := EarRepo.DeleteEar(uuid.NewV4().String())
		assert.Error(err, "Error was not thrown while deleting in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		err := EarRepo.DeleteEar(uuid)
		assert.Error(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		err := EarRepo.DeleteEar(uuid)
		assert.Error(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		EarRepo.AddEar(&testEar)
		err := EarRepo.DeleteEar(testEar.ID)
		assert.NoError(err, "Ear was not deleted.")
	})
}
