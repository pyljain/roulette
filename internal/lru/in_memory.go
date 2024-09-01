package lru

import (
	"errors"
	"slices"
)

type inMemory struct {
	maxSizeKb     float64
	items         []CacheItem
	currentSizeKb float64
}

func NewInMemory(maxSizeKb float64) *inMemory {
	return &inMemory{maxSizeKb: maxSizeKb}
}

var ErrNotExists error = errors.New("item does not exist")
var ErrTooLarge error = errors.New("item is too large")

func (cache *inMemory) Get(fileName string) ([]byte, error) {
	// Go through each item to check if exists
	for i, mi := range cache.items {
		if mi.Filename == fileName {
			// If exists then add to front
			cache.items = append(slices.Delete(cache.items, i, i+1), mi)
			return mi.Contents, nil
		}
	}

	// If not exists then return error
	return nil, ErrNotExists
}
func (im *inMemory) Add(filename string, contents []byte) error {
	// Check size of add item. If greater than total length then early return
	sizeInKB := float64(len(contents)) / 1024.0
	if sizeInKB > float64(im.maxSizeKb) {
		return ErrTooLarge
	}

	item := CacheItem{
		Filename: filename,
		Contents: contents,
		SizeInKB: sizeInKB,
	}
	im.items = append(im.items, item)
	im.currentSizeKb += sizeInKB
	// If not then compute total size of array and decide how to many to remove
	if im.currentSizeKb < float64(im.maxSizeKb) {
		return nil
	}

	// Trim array
	for im.currentSizeKb >= float64(im.maxSizeKb) {
		firstItem := im.items[0]
		im.items = slices.Delete(im.items, 0, 1)
		im.currentSizeKb -= firstItem.SizeInKB
	}

	return nil
}
