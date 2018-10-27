package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/SchiffFlieger/gophile/random"
	"github.com/pkg/errors"
)

var (
	ErrInvalidCredentials  = errors.New("Invalid Credentials")
	ErrEmptyUserOrPassword = errors.New("Empty Username or Password")
	ErrUserAlreadyExists   = errors.New("User already exists")
	GlobalAccountManager   Authenticator
)

const USERDIR = "accountsRoot"

type Authenticator interface {
	Authenticate(user, password string) bool
	CreateAccount(user, password string) error
	DeleteAccount(user, password string) error
	ChangePassword(user, oldPw, newPw string) error
	DeleteFile()
}

type AccountManager struct {
	lock       sync.RWMutex
	filename   string
	saltLength int
}

type Accounts struct {
	XMLName xml.Name  `xml:"accounts"`
	Accs    []Account `xml:"account"`
}

type Account struct {
	User     string `xml:"username"`
	Password string `xml:"password"`
	Salt     string `xml:"salt"`
}

// Creates a new account manager. Filename is the name of the file in which the accounts are stored.
// SaltLength determines the length of the salt for the passwords.
func NewAccountManager(filename string, saltLength int) Authenticator {
	return &AccountManager{filename: filename, saltLength: saltLength}
}

// Creates a new user account. Username and password must not be empty. A username has to
// be unique.
func (am *AccountManager) CreateAccount(user, password string) error {
	if user == "" || password == "" {
		return ErrEmptyUserOrPassword
	}

	accs := am.readAccounts()

	salt := random.Rnd.RandomString(am.saltLength)
	account := Account{User: user, Password: hashPassword(password, salt), Salt: salt}
	if contains(accs.Accs, account.User) {
		return ErrUserAlreadyExists
	}
	accs.Accs = append(accs.Accs, account)

	am.writeAccounts(accs)
	os.MkdirAll(path.Join(USERDIR, user), os.ModePerm)

	return nil
}

// Deletes a user account. Before deleting, this method checks if the given credentials are correct.
func (am *AccountManager) DeleteAccount(user, password string) error {
	if !am.Authenticate(user, password) {
		return ErrInvalidCredentials
	}

	accs := am.readAccounts()
	pos := getPosition(accs.Accs, user)
	accs.Accs = append(accs.Accs[:pos], accs.Accs[pos+1:]...)

	am.writeAccounts(accs)
	os.Remove(path.Join(USERDIR, user))

	return nil
}

// Changes the password of an existing user. The new password must not be empty.
func (am *AccountManager) ChangePassword(user, oldPw, newPw string) error {
	if newPw == "" {
		return ErrEmptyUserOrPassword
	}

	if !am.Authenticate(user, oldPw) {
		return ErrInvalidCredentials
	}

	accs := am.readAccounts()
	account := &accs.Accs[getPosition(accs.Accs, user)]
	account.Password = hashPassword(newPw, account.Salt)

	am.writeAccounts(accs)
	return nil
}

// Checks if the given credentials are correct.
func (am *AccountManager) Authenticate(user, password string) bool {
	accs := am.readAccounts()
	if !contains(accs.Accs, user) {
		return false
	}

	account := accs.Accs[getPosition(accs.Accs, user)]
	return hashPassword(password, account.Salt) == account.Password
}

// Deletes the files containing all the accounts.
func (am *AccountManager) DeleteFile() {
	am.lock.Lock()
	defer am.lock.Unlock()

	os.Remove(am.filename)
}

// Calculates a hash for a given password and salt.
func hashPassword(password string, salt string) string {
	hasher := sha256.New()
	io.WriteString(hasher, password)
	io.WriteString(hasher, salt)
	return hex.EncodeToString(hasher.Sum(nil))
}

// Returns the position of an account in the accounts list
func getPosition(accs []Account, user string) int {
	for p, v := range accs {
		if v.User == user {
			return p
		}
	}
	return -1
}

// Checks if the given username belongs to an existing user.
func contains(accs []Account, user string) bool {
	return getPosition(accs, user) >= 0
}

// Reads all accounts from the accounts file. If there are any errors while reading, this method
// returns an empty account list.
func (am *AccountManager) readAccounts() Accounts {
	am.lock.RLock()
	defer am.lock.RUnlock()

	file, err := os.Open(am.filename)
	if err != nil {
		return Accounts{}
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return Accounts{}
	}

	accounts := Accounts{}
	err = xml.Unmarshal(data, &accounts)
	if err != nil {
		return Accounts{}
	}

	return accounts
}

// Writes all existing accounts to the account file.
func (am *AccountManager) writeAccounts(accs Accounts) error {
	am.lock.Lock()
	defer am.lock.Unlock()

	output, err := xml.MarshalIndent(accs, "", "	")
	if err != nil {
		return errors.Wrap(err, "Error while marshalling")
	}

	file, err := os.Create(am.filename)
	if os.IsNotExist(err) {
		os.MkdirAll(USERDIR, os.ModePerm)
		file, err = os.Create(am.filename)
		if err != nil {
			return errors.Wrap(err, "Can not create accounts file")
		}
	}
	defer file.Close()

	output = append([]byte(xml.Header), output...)

	_, err = file.Write(output)
	if err != nil {
		return errors.Wrap(err, "Can not write accounts file")
	}

	return nil
}
