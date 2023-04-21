package entities_test

import (
	"testing"

	"github.com/LucasGois1/jarvis/src/domain/entities"
)

func TestNewRoleShouldReturnARespectiveRoleWhenGivenStringSystem(t *testing.T) {

	role, _ := entities.NewRole("System")

	if role != entities.System {
		t.Errorf("expected %v, got %v", entities.System, role)
	}
}
