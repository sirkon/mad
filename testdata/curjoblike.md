# Service.Method

### requests
* first one of variant with `a` field
```json
{
    "a": 1
}
```

* second one of variant with `b` field
```json
{
    "b": "field"
}
```

### response(OK)
```json
{
    "status": "OK",
    "fieldConsumed": "a"
}
```

### response(ERROR, NOT_AVALABLE)
Only one sample for status `ERROR` is given as others has the same scheme
```json
{
    "status": "ERROR",
    "reason": "DB error: incorrect query syntax at …"
}
```

