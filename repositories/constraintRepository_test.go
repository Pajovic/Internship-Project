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
		utils.DropTables(ConstraintRepo.DB)
		defer SetupTables(ConstraintRepo.DB)
		err := ConstraintRepo.AddConstraint(&testConstraint)
		assert.Error(err, "Error was not thrown while inserting in non-existing table")
	})

	t.Run("successful query", func(t *testing.T) {
		oldConstraints, _ := ConstraintRepo.GetAllConstraints()
		err := ConstraintRepo.AddConstraint(&testConstraint)
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
		utils.DropTables(ConstraintRepo.DB)
		defer SetupTables(ConstraintRepo.DB)
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
		ConstraintRepo.AddConstraint(&testConstraint)
		constraint, err := ConstraintRepo.GetConstraint(testConstraint.ID)
		assert.NotNil(constraint, "Result is nil")
		assert.NoError(err, "There was error while getting constraint")
		assert.Equal(testConstraint.ID, constraint.ID, "Returned constraint ID and test ID do not match.")
	})
}

func TestUpdateConstraint(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(ConstraintRepo.DB)
		defer SetupTables(ConstraintRepo.DB)
		err := ConstraintRepo.UpdateConstraint(testConstraint)
		assert.Error(err, "Error was not thrown while updating in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		testConstraint.ID = uuid
		err := ConstraintRepo.UpdateConstraint(testConstraint)
		assert.NotNil(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		testConstraint.ID = uuid
		err := ConstraintRepo.UpdateConstraint(testConstraint)
		assert.NotNil(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		ConstraintRepo.AddConstraint(&testConstraint)
		testConstraint.PropertyValue = 30
		err := ConstraintRepo.UpdateConstraint(testConstraint)
		assert.NoError(err, "Constraint was not updated.")
	})

}

func TestDeleteConstraint(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(ConstraintRepo.DB)
		defer SetupTables(ConstraintRepo.DB)
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
		ConstraintRepo.AddConstraint(&testConstraint)
		err := ConstraintRepo.DeleteConstraint(testConstraint.ID)
		assert.NoError(err, "Constraint was not deleted.")
	})
}
