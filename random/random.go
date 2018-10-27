package random

import (
	"math/rand"
	"sync"
	"time"
)

var Rnd = NewRandomizer()

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const charlen = 62

type Randomizer struct {
	lock   sync.Mutex
	random *rand.Rand
}

func NewRandomizer() *Randomizer {
	return &Randomizer{random: rand.New(rand.NewSource(time.Now().UTC().UnixNano()))}
}

// Creates a random string with given length. The string consists of uppercase letters, lowercase letters and digits.
func (rnd *Randomizer) RandomString(size int) string {
	rnd.lock.Lock()
	defer rnd.lock.Unlock()

	var result []byte

	for i := 0; i < size; i++ {
		char := chars[rnd.random.Int31n(charlen)]
		result = append(result, char)
	}

	return string(result)
}
