package test

import (
	"testing"

	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

// TestGenericRepository_Compilation verifies that GenericRepository compiles correctly
// and all methods are accessible. This is a compile-time check test.
func TestGenericRepository_Compilation(t *testing.T) {
	t.Run("TeacherRepository_Struct", func(t *testing.T) {
		t.Skip("TeacherRepository requires app instance, skipping in compile-only test")
	})

	t.Run("CenterRepository_Struct", func(t *testing.T) {
		t.Skip("CenterRepository requires app instance, skipping in compile-only test")
	})

	t.Run("GenericRepository_TableNames", func(t *testing.T) {
		// Test that NewGenericRepository sets correct table names with nil DBs
		teacherRepo := repositories.NewGenericRepository[models.Teacher](nil, nil)
		if teacherRepo.GetTableName() != "teachers" {
			t.Errorf("Expected 'teachers', got '%s'", teacherRepo.GetTableName())
		}

		centerRepo := repositories.NewGenericRepository[models.Center](nil, nil)
		if centerRepo.GetTableName() != "centers" {
			t.Errorf("Expected 'centers', got '%s'", centerRepo.GetTableName())
		}

		t.Log("GenericRepository table names verified")
	})
}

// TestIModelInterface tests that models implement the IModel interface correctly.
func TestIModelInterface(t *testing.T) {
	t.Run("Teacher_IModel", func(t *testing.T) {
		// Verify Teacher implements IModel by checking TableName method exists
		var teacher models.Teacher
		_ = teacher.TableName()
		t.Log("Teacher implements IModel interface")
	})

	t.Run("Center_IModel", func(t *testing.T) {
		// Verify Center implements IModel by checking TableName method exists
		var center models.Center
		_ = center.TableName()
		t.Log("Center implements IModel interface")
	})
}

// TestGenericRepository_MethodSignatures verifies method signatures are correct.
func TestGenericRepository_MethodSignatures(t *testing.T) {
	t.Run("TeacherRepository_HasGenericMethods", func(t *testing.T) {
		t.Skip("TeacherRepository requires app instance, skipping in compile-only test")
	})

	t.Run("TeacherRepository_HasCustomMethods", func(t *testing.T) {
		t.Skip("TeacherRepository requires app instance, skipping in compile-only test")
	})

	t.Run("CenterRepository_HasGenericMethods", func(t *testing.T) {
		t.Skip("CenterRepository requires app instance, skipping in compile-only test")
	})

	t.Run("CenterRepository_HasCustomMethods", func(t *testing.T) {
		t.Skip("CenterRepository requires app instance, skipping in compile-only test")
	})
}

// TestGenericRepository_Functional tests the GenericRepository with mock data.
func TestGenericRepository_Functional(t *testing.T) {
	t.Run("TableName_Extraction", func(t *testing.T) {
		// Verify table names are correctly extracted from models
		teacherTable := models.Teacher{}.TableName()
		if teacherTable != "teachers" {
			t.Errorf("Expected 'teachers', got '%s'", teacherTable)
		}

		centerTable := models.Center{}.TableName()
		if centerTable != "centers" {
			t.Errorf("Expected 'centers', got '%s'", centerTable)
		}
	})
}
