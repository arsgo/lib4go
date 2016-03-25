package memcache

import (
    "github.com/bradfitz/gomemcache/memcache"
)


type MemcacheClient struct {
    server string
    client *memcache.Client
}

func New(server string)(*MemcacheClient){
    m:=&MemcacheClient{server:server}
    m.client=memcache.New(server)  
    return m
}

func (c *MemcacheClient) Get(key string)(string) {   
    data, err := c.client.Get(key)
    if(err != nil) {     
        return ""
    }
    return string(data.Value)
}

func (c *MemcacheClient) Set(key string, value string, expiresAt int32)(error) {
    data := &memcache.Item{Key: key, Value: []byte(value), Expiration : expiresAt}   
    return c.client.Set(data)
}

func (c *MemcacheClient) Delete(key string)(error){   
   return c.client.Delete(key)
}

func (c *MemcacheClient) Delay(key string, expiresAt int32)(error) {
    v:=c.Get(key)
    data := &memcache.Item{Key: key, Value: []byte(v), Expiration : expiresAt}   
   return c.client.Set(data)
}



