# JSON vector

JSON parser based on [Vector API](https://github.com/koykov/vector) with minimum memory consumption.

### Usage

```go
src := []byte(`{"id":1,"name":"Foo","price":123,"tags":["Bar","Eek"],"stock":{"warehouse":300,"retail":20}}`)
vec := jsonvector.Acquire()
defer jsonvector.Release(vec)
_ = vec.Parse(src)
fmt.Println(vec.Dot("stock.warehouse")) // 300
```
