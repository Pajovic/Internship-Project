package repositories

import (
	"internship_project/models"
	"internship_project/utils"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddConstraint(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := ConstraintRepo.AddConstraint(&utils.TestConstraint)
		assert.Error(err, "Error was not thrown while inserting in non-existing table")
	})

	t.Run("successful query", func(t *testing.T) {
		oldConstraints, _ := ConstraintRepo.GetAllConstraints()
		err := ConstraintRepo.AddConstraint(&utils.TestConstraint)
		newConstraints, _ := ConstraintRepo.GetAllConstraints()

		assert.NoError(err)
		assert.Equal(1, len(newConstraints)-len(oldConstraints), "Constraint was not added.")

	})

}

func TestGetAllConstraints(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful query", func(t *testing.T) {
		allConstraints, err := ConstraintRepo.GetAllConstraints()

		assert.NoError(err, "Error was thrown while reading constraints")
		assert.NotNil(allConstraints, "Constraints returned are nil.")
		assert.IsType(allConstraints, []models.AccessConstraint{}, "Returned result is not of type Constraint")

	})
}

func TestGetConstraint(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		_, err := ConstraintRepo.GetConstraint(uuid.NewV4().String())
		assert.Error(err, "Error was not thrown while getting from non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		_, err := ConstraintRepo.GetConstraint(uuid)
		assert.Error(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		_, err := ConstraintRepo.GetConstraint(uuid)
		assert.Error(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		ConstraintRepo.AddConstraint(&utils.TestConstraint)
		constraint, err := ConstraintRepo.GetConstraint(utils.TestConstraint.ID)
		assert.NotNil(constraint, "Result is nil")
		assert.NoError(err, "There was error while getting constraint")
		assert.Equal(utils.TestConstraint.ID, constraint.ID, "Returned constraint ID and test ID do not match.")
	})
}

func TestUpdateConstraint(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := ConstraintRepo.UpdateConstraint(utils.TestConstraint)
		assert.Error(err, "Error was not thrown while updating in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		utils.TestConstraint.ID = uuid
		err := ConstraintRepo.UpdateConstraint(utils.TestConstraint)
		assert.NotNil(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		utils.TestConstraint.ID = uuid
		err := ConstraintRepo.UpdateConstraint(utils.TestConstraint)
		assert.NotNil(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		ConstraintRepo.AddConstraint(&utils.TestConstraint)
		utils.TestConstraint.PropertyValue = 30
		err := ConstraintRepo.UpdateConstraint(utils.TestConstraint)
		assert.NoError(err, "Constraint was not updated.")
	})

}

func TestDeleteConstraint(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := ConstraintRepo.DeleteConstraint(uuid.NewV4().String())
		assert.Error(err, "Error was not thrown while deleting in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		err := ConstraintRepo.DeleteConstraint(uuid)
		assert.Error(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		err := ConstraintRepo.DeleteConstraint(uuid)
		assert.Error(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		ConstraintRepo.AddConstraint(&utils.TestConstraint)
		err := ConstraintRepo.DeleteConstraint(utils.TestConstraint.ID)
		assert.NoError(err, "Constraint was not deleted.")
	})
}
