package test

import (
	"testing"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
)

// ============ Center Teacher Note (Rating) Repository Tests ============

func TestCenterTeacherNoteRepository_Initialization(t *testing.T) {
	// This test verifies the repository can be instantiated
	repo := repositories.NewCenterTeacherNoteRepository(nil)
	if repo == nil {
		t.Error("Expected repository to be created")
	}
}

func TestCenterTeacherNoteModel_Structure(t *testing.T) {
	// Test that the model has the expected fields
	note := models.CenterTeacherNote{
		CenterID:     1,
		TeacherID:    2,
		Rating:       5,
		InternalNote: "Test note",
	}

	if note.CenterID != 1 {
		t.Errorf("Expected CenterID 1, got %d", note.CenterID)
	}

	if note.TeacherID != 2 {
		t.Errorf("Expected TeacherID 2, got %d", note.TeacherID)
	}

	if note.Rating != 5 {
		t.Errorf("Expected Rating 5, got %d", note.Rating)
	}

	if note.InternalNote != "Test note" {
		t.Errorf("Expected InternalNote 'Test note', got '%s'", note.InternalNote)
	}
}

func TestCenterTeacherNote_TableName(t *testing.T) {
	note := models.CenterTeacherNote{}
	expectedTableName := "center_teacher_notes"

	if note.TableName() != expectedTableName {
		t.Errorf("Expected TableName '%s', got '%s'", expectedTableName, note.TableName())
	}
}

func TestCenterMembershipModel_Structure(t *testing.T) {
	membership := models.CenterMembership{
		ID:        1,
		CenterID:  2,
		TeacherID: 3,
		Status:    "ACTIVE",
	}

	if membership.ID != 1 {
		t.Errorf("Expected ID 1, got %d", membership.ID)
	}

	if membership.CenterID != 2 {
		t.Errorf("Expected CenterID 2, got %d", membership.CenterID)
	}

	if membership.TeacherID != 3 {
		t.Errorf("Expected TeacherID 3, got %d", membership.TeacherID)
	}

	if membership.Status != "ACTIVE" {
		t.Errorf("Expected Status 'ACTIVE', got '%s'", membership.Status)
	}
}

func TestCenterMembership_TableName(t *testing.T) {
	membership := models.CenterMembership{}
	expectedTableName := "center_memberships"

	if membership.TableName() != expectedTableName {
		t.Errorf("Expected TableName '%s', got '%s'", expectedTableName, membership.TableName())
	}
}

func TestErrInfos_ErrorCodes(t *testing.T) {
	// Test that error codes are defined correctly
	e := errInfos.Initialize(1)

	// Verify the error info system is initialized
	if e == nil {
		t.Error("Expected errInfos to be initialized")
	}
}
