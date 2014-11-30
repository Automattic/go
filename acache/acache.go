// Package cache is a simple wrapper for accessing gomemcache
// The package provides gob encoding/decoding of objects to make it easier
// to store structs in memcache
package acache

import (
	"bytes"
	"encoding/gob"

	"github.com/bradfitz/gomemcache/memcache"
)

type Cache struct {
	client *memcache.Client
}

func (c *Cache) Get(key string, v interface{}) (err error) {
	if c.client == nil {
		return memcache.ErrNoServers
	}

	item, err := c.client.Get(key)
	if err != nil {
		return err
	}
	decBuf := bytes.NewBuffer(item.Value)
	err = gob.NewDecoder(decBuf).Decode(v)
	if err != nil {
		// Failed to decode
		return err
	}
	return nil
}

func (c *Cache) Set(key string, data interface{}) (err error) {
	if c.client == nil {
		return memcache.ErrNoServers
	}

	encBuf := new(bytes.Buffer)
	err = gob.NewEncoder(encBuf).Encode(data)
	if err != nil {
		// Failed to encode value
		return err
	}

	item := memcache.Item{
		Key:        key,
		Value:      encBuf.Bytes(),
		Expiration: 60 * 60, // expiration in seconds
	}
	c.client.Set(&item)
	return err
}

func (c *Cache) Delete(key string) (err error) {
	if c.client == nil {
		return memcache.ErrNoServers
	}

	return c.client.Delete(key)
}

func (c *Cache) Connect(servers ...string) (err error) {
	c.client = memcache.New(servers...)
	return err
}
