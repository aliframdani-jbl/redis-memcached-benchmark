In default it tested get set operation in Memcached and Redis String. If you want to switch the test to Redis Sets, rename the corresponding file to `*_test.go` and remove the `_test` suffix filename on the untested case.

Ex:
```
redis_sets.go -> redis_sets_test.go
redis_string_test.go - > redis_string.go
```