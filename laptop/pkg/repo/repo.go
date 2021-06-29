package repo

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"

	"github.com/chutommy/simple-bank/laptop/laptop"
)

var (
	ErrAlreadyExists = errors.New("record already exists")
)

type Repo interface {
	CreateLaptop(l *laptop.Laptop) error
}

type inMemoryRepo struct {
	mutex   sync.Mutex
	laptops map[string]*laptop.Laptop
}

func NewRepo() Repo {
	return &inMemoryRepo{
		laptops: make(map[string]*laptop.Laptop),
	}
}

func (r *inMemoryRepo) CreateLaptop(l *laptop.Laptop) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// check if unique id
	if _, ok := r.laptops[l.Id]; ok {
		return fmt.Errorf("%w: id: %s", ErrAlreadyExists, l.Id)
	}

	var l2 *laptop.Laptop
	err := copier.Copy(l2, l)
	if err != nil {
		return fmt.Errorf("cannot create a storage copy of the laptop: %w", err)
	}

	r.laptops[l.Id] = l

	return nil
}
