# LocalLRUCache - an in-memory LRU cache with expiration and length for any data type


### usage

``` 
func main() {
    if cache, ok := InitLRUCache(size, timeout); ok { // size > 0;  timeout=0 no timeout
		cache.Set("key","value")
		value, found:=cache.Get("key")
		value, found:=cache.Peek("key")
		count := cache.Count()
		cache.Remove("key")
		cache.Clear()
	}
}
``` 
