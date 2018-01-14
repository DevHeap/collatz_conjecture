package cache

import (
	"sort"
)

type Index struct {
	payload []record
}


type record struct {
	key    string
	length int
}

func NewIndex() Index{
	return Index{payload: make([]record, 0)}
}


func (i *Index) Update(key string, length int) {
	element := record{key:key, length: length}

	// search for proper position of element in index
	position := sort.Search(len(i.payload), func(j int) bool { return i.payload[j].length > element.length })
	// extend index
	i.payload = append(i.payload, record{})
	// shift elements by copy
	copy(i.payload[position+1:], i.payload[position:])
	// set new element on proper position
	i.payload[position] = element
}

func (i *Index) Size() int {
	return len(i.payload)
}

func (i *Index) MinKey() string {
	return i.payload[0].key
}

func (i *Index) MinLength() int {
	return  i.payload[0].length
}

func (i *Index) RemMin() {
	if len(i.payload) > 0 {
		i.payload = append(i.payload[:0], i.payload[1:]...)
	}
}