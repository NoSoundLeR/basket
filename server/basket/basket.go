package basket

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
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
	Baskets []*Basket
}

// GetBasket ...
func (bs *Baskets) GetBasket(id string) (*Basket, error) {
	for _, v := range bs.Baskets {
		if v.ID.String() == id {
			return v, nil
		}
	}
	return nil, errors.New("could not find element with this ID")
}

// AddBasket ...
func (bs *Baskets) AddBasket(b Basket) {
	bs.Baskets = append(bs.Baskets, &b)
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
	ID     uuid.UUID
	Secret string

	// ProtectionLevel=0 - none
	// ProtectionLevel=1 - ip
	// ProtectionLevel=2 - cookie
	// ProtectionLevel=3 - ip + cookie
	ProtectionLevel int
	Active          bool
	Title           string
	Description     string
	CreatedAt       int64
	Timeout         int
	ClosedAt        int64
	Vars            map[string]string
	VarsZero        []string
	Result          string
}

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const randomStringLen = 16

// Data ...
type Data struct {
	ProtectionLevel string `json:"protectionLevel"`
	Timeout         string `json:"timeout"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Value           string `json:"value"`
}

func randomString() string {
	b := make([]byte, randomStringLen)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

// NewBasket ...
func NewBasket(ip string, d Data) (*Basket, error) {
	if err := checkLength(d.Value); err != nil {
		return nil, err
	}
	protectionLevel, err := strconv.Atoi(d.ProtectionLevel)
	if err != nil {
		return nil, errors.New("wrong protection level")
	}
	if err := checkProtectionLevel(protectionLevel); err != nil {
		return nil, err
	}
	timeout, err := strconv.Atoi(d.Timeout)
	if err != nil {
		return nil, errors.New("wrong timeout value")
	}
	if err := checkTimeout(timeout); err != nil {
		return nil, err
	}

	s := make([]string, 0, 2)
	s = append(s, d.Value)
	now := time.Now().Unix()
	return &Basket{
		ID:              uuid.New(),
		Secret:          randomString(),
		ProtectionLevel: protectionLevel,
		Active:          true,
		Title:           d.Title,
		Description:     d.Description,
		CreatedAt:       now,
		Timeout:         timeout,
		ClosedAt:        0,
		Vars:            map[string]string{ip: d.Value},
		VarsZero:        s,
		Result:          "",
	}, nil
}

// Add ...
func (b *Basket) Add(ip string, v string) bool {
	if b.Active {
		if b.ProtectionLevel == 0 {
			b.VarsZero = append(b.VarsZero, v)
		} else {
			b.Vars[ip] = v
		}
		return true
	}
	return false
}

// GetSize ...
func (b *Basket) GetSize() int {
	if b.ProtectionLevel == 0 {
		return len(b.VarsZero)
	}
	return len(b.Vars)
}

// getRandom ...
func (b *Basket) getRandom() (string, error) {
	rand.Seed(time.Now().Unix())
	if b.ProtectionLevel == 0 {
		return b.VarsZero[rand.Intn(len(b.VarsZero))], nil
	}
	// accessing random elements of maps is weird
	// so we are making a slice of values
	buf := make([]string, len(b.Vars))
	var i int
	for _, value := range b.Vars {
		buf[i] = value
		i++
	}
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

func checkProtectionLevel(l int) error {
	if l < 0 || l > 3 {
		return fmt.Errorf("invalid protection level")
	}
	return nil
}

func checkTimeout(t int) error {
	if t < 0 || t > 60*15 {
		return fmt.Errorf("timeout: 0 - 15 min")
	}
	return nil
}
