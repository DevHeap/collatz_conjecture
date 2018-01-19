package cache

import (
	"math/big"
	"log"
)

type Storage struct {
	payload map[string]*data
}

type data struct {
	count int
	path  []*big.Int
}

func NewStorage() Storage {
	return Storage{payload: make(map[string]*data)}
}

func (s *Storage) Put(path []*big.Int){
	for index, element := range path {
		s.put(element.String(), path[index:])
	}
}

func (s *Storage) put(key string, value []*big.Int){
	_, exist := s.payload[key]

	if exist {
		s.payload[key].count += 1
		s.payload[key].path = value

	} else {
		s.payload[key] = &data{count: 0, path: value}
	}
}

func (s *Storage) Rem(key string) {
	d, exist := s.payload[key]

	if !exist {
		return
	}

	for _, element := range d.path {
		s.rem(element.String())
	}
}

func (s *Storage) rem(key string){
	d, exist := s.payload[key]

	if !exist {
		return
	}

	d.count -= 1

	if d.count < 1 {
		delete(s.payload, key)
	}
}

func (s *Storage) Get(number *big.Int) ([]*big.Int, bool){
	key := number.String()

	if d, ok := s.payload[key]; ok {
		log.Println("Cache hit" + key)
		return d.path, true
	}
	return nil, false
}