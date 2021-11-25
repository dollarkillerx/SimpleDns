# SimpleDns Api Doc

## list

- [auth](auth)
- [获取用户id](获取用户id)
- [获取所有DNS规则](获取所有DNS规则)
- [删除DNS规则](删除DNS规则)
- [更新DNS规则](更新DNS规则)
- [添加DNS规则](添加DNS规则)

### 注意点

请求 返回状态码 不为200 返回结构皆为 String

### auth

GET `/auth`
Header: `token: string`

返回200 为验证成功  其他返回值 均为错误

### 获取用户id

GET `/ip`

200 成功 返回 用户IP

### 获取所有DNS规则   (需AUTH)

GET `/api/list`

Response:
``` 
[]model.DnsModel   的JSON


type DnsModel struct {
	ID        string        `json:"id"`
	Domain    string        `json:"domain"`
	QueryType easy_dns.Type `json:"query_type"`
	Addr      string        `json:"addr"`
	TTL       int           `json:"ttl"`
}
```


### 删除DNS规则 (需AUTH)

POST `/api/del`

Request:
```json
{
  "id": "id"    
}
```

Response:

ctx.String(200, "success")

### 更新DNS规则

POST `/api/update`

Request:
```json
type DnsModel struct {
  ID        string        `json:"id"`
  Domain    string        `json:"domain"`
  QueryType easy_dns.Type `json:"query_type"`
  Addr      string        `json:"addr"`
  TTL       int           `json:"ttl"`
}
```

Response:

ctx.String(200, "success")

### 添加DNS规则

POST `/api/add`

Request:
```json
type DnsModel struct {
  Domain    string        `json:"domain"`
  QueryType easy_dns.Type `json:"query_type"`
  Addr      string        `json:"addr"`
  TTL       int           `json:"ttl"`
}
```

Response:

ctx.String(200, "success")
