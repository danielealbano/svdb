package collection

import (
	"fmt"
	"github.com/c2h5oh/datasize"
	shared_support "github.com/danielealbano/cvdb/shared/support"
	usearch "github.com/unum-cloud/usearch/golang"
	"unsafe"
)

type Key usearch.Key

type Vector []float32

type Collection struct {
	index  *usearch.Index
	isFull bool
	Config *CollectionConfig
}

func NewCollection(config *CollectionConfig) (*Collection, error) {
	index, err := usearch.NewIndex(config.toUsearchConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return &Collection{
		index:  index,
		Config: config,
	}, nil
}

func (c *Collection) Load(path string) error {
	var size uint

	err := c.index.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load collection from path: %w", err)
	}

	// Get the current size of the index
	size, err = c.index.SerializedLength()
	if err != nil {
		return fmt.Errorf("failed to get size of index: %w", err)
	}

	if size >= c.Config.MaxSize {
		c.isFull = true
	}

	return nil
}

func (c *Collection) Destroy() error {
	if c.index == nil {
		return fmt.Errorf("collection not initialized")
	}

	err := c.index.Destroy()
	if err != nil {
		return fmt.Errorf("failed to destroy collection: %w", err)
	}

	return nil
}

func (c *Collection) Search(query Vector, limit uint32) ([]Key, []float32, error) {
	keys, distances, err := c.index.Search(query, uint(limit))

	if err != nil {
		return nil, nil, fmt.Errorf("failed to search: %w", err)
	}

	return *(*[]Key)(unsafe.Pointer(&keys)), distances, nil
}

func (c *Collection) Add(key Key, vector Vector) (uint64, bool, error) {
	return c.AddMulti([]Key{key}, []Vector{vector})
}

func (c *Collection) AddMulti(keys []Key, vectors []Vector) (uint64, bool, error) {
	var err error
	var initialSize uint
	var finalSize uint
	inserted := uint64(0)

	// TODO: The mechanism is not efficient at all, if 10000 vectors are added and the reservation triggers a growth
	//       of the index, only the first vector will be written and the rest will be skipped.
	//       To avoid wasting too much space, the code that follows, if the max size hasn't been reached, will get the
	//       current size (initialSize) and then check if the max size threshold has been reached only if the size
	//       changes again.
	//       To be efficient it's necessary to calculate in advance the numbers of vectors that can be added before
	//       reaching the max size (can be done checking capacity) and then know if the next upsize will hit the max
	//       size, and if so, prevent the addition of the vectors.
	//       Will be improved in the future.

	err = c.index.Reserve(uint(len(keys)))
	if err != nil {
		return 0, false, fmt.Errorf("failed to reserve space in index: %w", err)
	}

	if c.isFull {
		return 0, true, nil
	}

	// Get the current size of the index
	initialSize, err = c.index.SerializedLength()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get size of index: %w", err)
	}

	for i, key := range keys {
		err = c.index.Add(usearch.Key(key), vectors[i])
		if err != nil {
			return inserted, false, fmt.Errorf("failed to add vector to index: %w", err)
		}

		inserted++

		finalSize, err = c.index.SerializedLength()
		if err != nil {
			return inserted, false, fmt.Errorf("failed to get size of index: %w", err)
		}
		if initialSize != finalSize && finalSize >= c.Config.MaxSize {
			c.isFull = true
			break
		}
	}

	finalSize, _ = c.index.SerializedLength()
	shared_support.Logger().Info().Msgf(
		"inserted %d key(s) to index, current size: %s, the shard is full: %t",
		inserted,
		datasize.ByteSize(finalSize).HumanReadable(),
		c.isFull)

	return inserted, c.isFull, nil
}

func (c *Collection) Get(key Key, count uint) (Vector, error) {
	vector, err := c.index.Get(usearch.Key(key), count)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector from index: %w", err)
	}

	return vector, nil
}

func (c *Collection) Has(key Key) bool {
	_, err := c.index.Get(usearch.Key(key), 1)
	if err != nil {
		return false
	}

	return true
}

func (c *Collection) Delete(key Key) error {
	err := c.index.Remove(usearch.Key(key))
	if err != nil {
		return fmt.Errorf("failed to delete vector from index: %w", err)
	}

	return nil
}

func (c *Collection) Length() (uint, error) {
	length, err := c.index.Len()
	if err != nil {
		return 0, fmt.Errorf("failed to get length of index: %w", err)
	}

	return length, nil
}

func (c *Collection) Capacity() (uint, error) {
	capacity, err := c.index.Capacity()
	if err != nil {
		return 0, fmt.Errorf("failed to get capacity of index: %w", err)
	}

	return capacity, nil
}

func (c *Collection) Size() (uint, error) {
	size, err := c.index.SerializedLength()
	if err != nil {
		return 0, fmt.Errorf("failed to get size of index: %w", err)
	}

	return size, nil
}

func (c *Collection) Save(path string) error {
	err := c.index.Save(path)

	if err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	return nil
}
