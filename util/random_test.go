package util_test

import (
	"github.com/chutommy/simple-bank/util"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestRandomAmount(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	amount1 := util.RandomAmount()

	rand.Seed(time.Now().UnixNano())
	amount2 := util.RandomAmount()

	assert.NotEqual(t, amount1, amount2)
}

func TestRandomBalance(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	balance1 := util.RandomBalance()

	rand.Seed(time.Now().UnixNano())
	balance2 := util.RandomBalance()

	assert.NotEqual(t, balance1, balance2)
}

func TestRandomCurrency(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	currency1 := util.RandomCurrency()

	rand.Seed(time.Now().UnixNano())
	currency2 := util.RandomCurrency()

	assert.NotEqual(t, currency1, currency2)
}

func TestRandomEmail(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	email1 := util.RandomEmail()

	rand.Seed(time.Now().UnixNano())
	email2 := util.RandomEmail()

	assert.NotEqual(t, email1, email2)
}

func TestRandomInt(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	int1 := util.RandomInt(1, 4096)

	rand.Seed(time.Now().UnixNano())
	int2 := util.RandomInt(1, 4096)

	assert.NotEqual(t, int1, int2)
}

func TestRandomOwner(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	owner1 := util.RandomOwner()

	rand.Seed(time.Now().UnixNano())
	owner2 := util.RandomOwner()

	assert.NotEqual(t, owner1, owner2)
}
