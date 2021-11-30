![](./README/logo.png)

# SimpleDns

Simple DNS relay server with customizable routing tables
![](./README/SimpleDns.jpg)

web: https://github.com/XC-Zero/simple_dns

![](./README/img.png)

### deploy
``` 
docker-compose up -d
```


### routing table exp:
routing_table.csv
``` 
domain,ip
domain,ip
```

### runing
`go build cmd/simple_dns/simple_dns.go`


### 关于ubuntu 端口占用

``` 
systemctl stop systemd-resolved.service
systemctl disable systemd-resolved.service
rm /etc/resolv.conf
echo "nameserver 8.8.8.8" >> /etc/resolv.conf
```
