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

const (
    // ResourceHeader.Type and Question.Type
    TypeA     Type = 1
    TypeNS    Type = 2
    TypeCNAME Type = 5
    TypeSOA   Type = 6
    TypePTR   Type = 12
    TypeMX    Type = 15
    TypeTXT   Type = 16
    TypeAAAA  Type = 28
    TypeSRV   Type = 33
    TypeOPT   Type = 41
    
    // Question.Type
    TypeWKS   Type = 11
    TypeHINFO Type = 13
    TypeMINFO Type = 14
    TypeAXFR  Type = 252
    TypeALL   Type = 255
)

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
