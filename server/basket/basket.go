package basket

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// Getter ...
type Getter interface {
	GetBasket(string) (*Basket, error)
}

// Creator ...
type Creator interface {
	AddBasket(Basket)
}

// Baskets ...
type Baskets struct {
	Baskets []Basket
}

// GetBasket ...
func (bs *Baskets) GetBasket(id string) (*Basket, error) {
	for _, v := range bs.Baskets {
		if v.ID.String() == id {
			return &v, nil
		}
	}
	return nil, errors.New("could not find element with this ID")
}

// AddBasket ...
func (bs *Baskets) AddBasket(b Basket) {
	bs.Baskets = append(bs.Baskets, b)
}

// RemoveBasket ...
func (bs *Baskets) RemoveBasket(id string) {
	idx := -1
	for i, v := range bs.Baskets {
		if v.ID.String() == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		panic("could not find element with this ID")
	}
	(bs.Baskets)[idx] = (bs.Baskets)[len(bs.Baskets)-1]
	bs.Baskets = (bs.Baskets)[:len(bs.Baskets)-1]
}

// Basket ...
type Basket struct {
	ID          uuid.UUID
	Secret      string
	Active      bool
	Title       string
	Description string
	CreatedAt   int64
	ClosedAt    int64
	Vars        map[string]string
	Result      string
}

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const randomStringLen = 16

func randomString() string {
	b := make([]byte, randomStringLen)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

// NewBasket ...
func NewBasket(ip string, title string, description string, v string) (*Basket, error) {
	if err := checkLength(v); err != nil {
		return nil, err
	}
	return &Basket{
		ID:          uuid.New(),
		Secret:      randomString(),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now().Unix(),
		ClosedAt:    0,
		Active:      true,
		Vars:        map[string]string{ip: v},
		Result:      "",
	}, nil
}

// Add ...
func (b *Basket) Add(ip string, v string) bool {
	if b.Active {
		b.Vars[ip] = v
		return true
	}
	return false
}

// getRandom ...
func (b *Basket) getRandom() (string, error) {
	if len(b.Vars) <= 1 {
		return "", errors.New("not enough vars")
	}
	// accessing random elements of maps is weird
	// so we are making a slice of values
	buf := make([]string, len(b.Vars))
	var i int
	for _, value := range b.Vars {
		buf[i] = value
		i++
	}
	rand.Seed(time.Now().Unix())
	return buf[rand.Intn(len(buf))], nil
}

// Close ...
func (b *Basket) Close() {
	result, err := b.getRandom()
	if err == nil {
		b.Result = result
	} else {
		b.Result = err.Error()
	}
	b.Active = false
	b.ClosedAt = time.Now().Unix()
}

const maxLen = 1024

func checkLength(s string) error {
	if len(s) > maxLen {
		return fmt.Errorf("length can't be more than %d", maxLen)
	}
	return nil
}
