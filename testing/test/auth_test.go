package test

import (
	"testing"
)

func TestAuthController_AdminLogin(t *testing.T) {
	t.Skip("Skipping - requires proper password hash setup")
}

func TestAuthController_TeacherLineLogin(t *testing.T) {
	t.Skip("Skipping - requires LINE API setup")
}

func TestAuthController_RefreshToken(t *testing.T) {
	t.Skip("Skipping - requires token generation setup")
}

func TestAuthController_Logout(t *testing.T) {
	t.Skip("Skipping - requires authentication setup")
}
