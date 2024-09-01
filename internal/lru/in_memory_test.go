package lru

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {

	tt := []struct {
		Description         string
		Filename            string
		Contents            []byte
		CacheSize           float64
		ExpectedError       bool
		ExpectedCacheLength int
		ExistingCachedItems []CacheItem
	}{
		{
			Description:         "Test adding a small file",
			Filename:            "./test.txt",
			Contents:            []byte("Hello"),
			CacheSize:           10,
			ExpectedError:       false,
			ExpectedCacheLength: 1,
			ExistingCachedItems: []CacheItem{},
		},
		{
			Description:         "Test adding a file larger than cache size",
			Filename:            "./test.txt",
			Contents:            []byte("The user proxy agent is used for interacting with the assistant agent"),
			CacheSize:           0.00001,
			ExpectedError:       true,
			ExpectedCacheLength: 0,
			ExistingCachedItems: []CacheItem{},
		},
		{
			Description:         "Test adding a file when the cache has existing items but is not full",
			Filename:            "./test.txt",
			Contents:            []byte("The user proxy agent is used for interacting with the assistant agent"),
			CacheSize:           0.1,
			ExpectedError:       false,
			ExpectedCacheLength: 3,
			ExistingCachedItems: []CacheItem{
				{
					Filename: "/pkg/text.json",
					Contents: []byte("Test"),
					SizeInKB: 0.02,
				},
				{
					Filename: "/pkg/text.txt",
					Contents: []byte("Test2"),
					SizeInKB: 0.01,
				},
			},
		},
		{
			Description:         "Test adding a file when the cache has existing items and is full, should evict item",
			Filename:            "./test.txt",
			Contents:            []byte("The user proxy agent is used for interacting with the assistant agent"),
			CacheSize:           0.1,
			ExpectedError:       false,
			ExpectedCacheLength: 2,
			ExistingCachedItems: []CacheItem{
				{
					Filename: "/pkg/text.json",
					Contents: []byte("Test"),
					SizeInKB: 0.05,
				},
				{
					Filename: "/pkg/text.txt",
					Contents: []byte("Test2"),
					SizeInKB: 0.03,
				},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.Description, func(t *testing.T) {
			cache := NewInMemory(test.CacheSize)
			cache.items = test.ExistingCachedItems
			cache.currentSizeKb = findSize(cache.items)
			err := cache.Add(test.Filename, test.Contents)
			l := len(cache.items)
			if test.ExpectedError {
				require.Error(t, err)
				require.Equal(t, test.ExpectedCacheLength, l)
				return
			}
			require.Equal(t, test.ExpectedCacheLength, l)
			require.NoError(t, err)
		})
	}
}

func findSize(items []CacheItem) float64 {
	totalSize := float64(0.0)
	for _, item := range items {
		totalSize += item.SizeInKB
	}

	return totalSize
}

func TestGet(t *testing.T) {
	tt := []struct {
		description            string
		fileName               string
		expectedError          bool
		expectedContents       []byte
		existingCacheFileNames []CacheItem
	}{
		{
			description:            "When file not in cache should throw error",
			fileName:               "./test.txt",
			expectedError:          true,
			expectedContents:       nil,
			existingCacheFileNames: []CacheItem{},
		},
		{
			description:      "When file in cache should return file contents ",
			fileName:         "./test.txt",
			expectedError:    false,
			expectedContents: []byte("Hello"),
			existingCacheFileNames: []CacheItem{
				{
					Filename: "./test.txt",
					Contents: []byte("Hello"),
					SizeInKB: 1,
				},
			},
		},
		{
			description:      "When file in cache and is not the latest should become the latest item",
			fileName:         "./one.txt",
			expectedError:    false,
			expectedContents: []byte("Hello"),
			existingCacheFileNames: []CacheItem{
				{
					Filename: "./one.txt",
					Contents: []byte("Hello"),
					SizeInKB: 1,
				},
				{
					Filename: "./two.txt",
					Contents: []byte("Hello"),
					SizeInKB: 1,
				},
				{
					Filename: "./three.txt",
					Contents: []byte("Hello"),
					SizeInKB: 1,
				},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			cache := NewInMemory(5)
			cache.items = test.existingCacheFileNames
			contents, err := cache.Get(test.fileName)
			if test.expectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedContents, contents)
			require.Equal(t, test.fileName, cache.items[len(cache.items)-1].Filename)
		})
	}
}
