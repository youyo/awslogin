package awslogin

import (
	"time"

	"github.com/adelowo/onecache"
	"github.com/adelowo/onecache/filesystem"
	"github.com/aws/aws-sdk-go/aws/credentials"
	homedir "github.com/mitchellh/go-homedir"
)

// Cache struct
type Cache struct {
	Store   *filesystem.FSStore
	Marshal *onecache.CacheSerializer
	Key     string
}

// NewCache
func NewCache(path, key string) (c *Cache, err error) {
	fullPath, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	store := filesystem.MustNewFSStore(fullPath)
	marshal := onecache.NewCacheSerializer()

	c = &Cache{
		Store:   store,
		Marshal: marshal,
		Key:     key,
	}

	return c, nil
}

// Save
func (c *Cache) Save(creds *credentials.Value) (err error) {
	dataByte, err := c.Marshal.Serialize(&creds)
	if err != nil {
		return err
	}

	err = c.Store.Set(c.Key, dataByte, 12*time.Hour)
	return err
}

// Load
func (c *Cache) Load() (creds *credentials.Value, err error) {
	credsByte, err := c.Store.Get(c.Key)
	if err != nil {
		return nil, err
	}

	c.Marshal.DeSerialize(credsByte, &creds)

	return creds, nil
}
