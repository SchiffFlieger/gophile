package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunAllTests(t *testing.T) {
	GlobalAccountManager = NewAccountManager("accounts.xml", 16)
	defer func() {
		GlobalAccountManager.DeleteFile()
		os.Remove(USERDIR)
	}()

	t.Run("add and delete account", addAndDeleteAccount)
	t.Run("try delete with invalid credentials", tryDeleteWithInvalidCredentials)
	t.Run("add account twice", addAccountTwice)
	t.Run("verify invalid account", verifyAccountInvalid)
	t.Run("verify account not exist", verifyAccountNotExist)
	t.Run("verify valid account", verifyAccountValid)
	t.Run("change user passwort with invalid credentials", changeUserPasswordInvalid)
	t.Run("change user passwort with valid credentials", changeUserPasswordValid)
	t.Run("change user passwort to empty password", changeUserPasswordEmpty)
	t.Run("empty username", emptyUsername)
	t.Run("empty password", emptyPassword)
}

func addAndDeleteAccount(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "123456")

	result := GlobalAccountManager.Authenticate("testuser", "123456")
	assert.Equal(t, true, result)
}

func tryDeleteWithInvalidCredentials(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "123456")

	GlobalAccountManager.DeleteAccount("testuser", "abcdef")
	result := GlobalAccountManager.Authenticate("testuser", "123456")
	assert.Equal(t, true, result)
}

func addAccountTwice(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "123456")

	err := GlobalAccountManager.CreateAccount("testuser", "123456")
	assert.Equal(t, ErrUserAlreadyExists, err)
}

func verifyAccountInvalid(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "123456")

	result := GlobalAccountManager.Authenticate("testuser", "654321")
	assert.Equal(t, false, result)
}

func verifyAccountNotExist(t *testing.T) {
	result := GlobalAccountManager.Authenticate("testuser", "123456")
	assert.Equal(t, false, result)
}

func verifyAccountValid(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "123456")

	result := GlobalAccountManager.Authenticate("testuser", "123456")
	assert.Equal(t, true, result)
}

func changeUserPasswordInvalid(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "123456")

	err := GlobalAccountManager.ChangePassword("testuser", "abcdef", "654321")
	assert.Equal(t, ErrInvalidCredentials, err)
}

func changeUserPasswordValid(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "abcdef")

	GlobalAccountManager.ChangePassword("testuser", "123456", "abcdef")
	result := GlobalAccountManager.Authenticate("testuser", "abcdef")
	assert.Equal(t, true, result)
}

func changeUserPasswordEmpty(t *testing.T) {
	GlobalAccountManager.CreateAccount("testuser", "123456")
	defer GlobalAccountManager.DeleteAccount("testuser", "123456")

	err := GlobalAccountManager.ChangePassword("testuser", "123456", "")
	assert.Equal(t, ErrEmptyUserOrPassword, err)

	result := GlobalAccountManager.Authenticate("testuser", "123456")
	assert.Equal(t, true, result)
}

func emptyUsername(t *testing.T) {
	err := GlobalAccountManager.CreateAccount("", "123456")
	assert.Equal(t, ErrEmptyUserOrPassword, err)
}

func emptyPassword(t *testing.T) {
	err := GlobalAccountManager.CreateAccount("testuser", "")
	assert.Equal(t, ErrEmptyUserOrPassword, err)
}
