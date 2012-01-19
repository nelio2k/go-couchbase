package couchbase

import (
	"encoding/json"
	"log"
)

func (b *Bucket) do(k string, f func(mc *memcachedClient) error) error {
	h := b.VBHash(k) % uint32(len(b.connections))
	log.Printf("Sending %s to connection %d (%d connections)\n", k, h, len(b.connections))
	if b.connections[h] == nil {
		b.connections[h] = connect("tcp", b.VBucketServerMap.ServerList[h])
	}
	return f(b.connections[h])
}

func (b *Bucket) Set(k string, v interface{}) error {
	return b.do(k, func(mc *memcachedClient) error {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		res := mc.Set(k, 0, 0, data)
		if res.Status != SUCCESS {
			return res
		}
		return nil
	})
}

func (b *Bucket) Get(k string, rv interface{}) error {
	return b.do(k, func(mc *memcachedClient) error {
		res := mc.Get(k)
		if res.Status != SUCCESS {
			return res
		}
		return json.Unmarshal(res.Body, rv)
	})
}

func (b *Bucket) Delete(k string) error {
	return b.do(k, func(mc *memcachedClient) error {
		res := mc.Del(k)
		if res.Status != SUCCESS {
			return res
		}
		return nil
	})
}
