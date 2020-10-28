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
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := EarRepo.AddEar(&utils.TestEar)
		assert.Error(err, "Error was not thrown while inserting in non-existing table")
	})

	t.Run("successful query", func(t *testing.T) {
		oldEars, _ := EarRepo.GetAllEars()
		err := EarRepo.AddEar(&utils.TestEar)
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
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
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
		EarRepo.AddEar(&utils.TestEar)
		ear, err := EarRepo.GetEar(utils.TestEar.ID)
		assert.NotNil(ear, "Result is nil")
		assert.NoError(err, "There was error while getting ear")
		assert.Equal(utils.TestEar.ID, ear.ID, "Returned ear ID and test ID do not match.")
	})
}

func TestUpdateEar(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := EarRepo.UpdateEar(utils.TestEar)
		assert.Error(err, "Error was not thrown while updating in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		utils.TestEar.ID = uuid
		err := EarRepo.UpdateEar(utils.TestEar)
		assert.NotNil(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		utils.TestEar.ID = uuid
		err := EarRepo.UpdateEar(utils.TestEar)
		assert.NotNil(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		EarRepo.AddEar(&utils.TestEar)
		utils.TestEar.Delete = true
		err := EarRepo.UpdateEar(utils.TestEar)
		assert.NoError(err, "Ear was not updated.")
	})

}

func TestDeleteEar(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
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
		EarRepo.AddEar(&utils.TestEar)
		err := EarRepo.DeleteEar(utils.TestEar.ID)
		assert.NoError(err, "Ear was not deleted.")
	})
}
