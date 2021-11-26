#!/bin/bash

source_path=./cmd/simple_dns
go_file=simple_dns.go
image_name=simple_dns:latest

package() {
  CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -o $source_path/simple_dns $source_path/$go_file
  docker rmi -f $image_name
  docker build -f $source_path/Dockerfile -t $image_name  .
  docker save -o simple_dns.tar $image_name

  rm -f $source_path/simple_dns
}

deploy() {
  docker load -i simple_dns.tar
}

# 开始菜单
start_menu() {
    echo "======================================"
    echo "Simple Dns     Build  by: DollarKiller"
    echo "======================================"
    echo "1.当前系统docker打包"
    echo "2.当前系统docker部署"
    echo "0.退出脚本"
    read -p "请输入数字" num
    case "$num" in
        1)
          echo "启动Docker打包"
          package
        ;;
        2)
          echo "Docker部署"
          deploy
        ;;
        0)
          exit 1
        ;;
        *)
        clear
        echo "请输入正确数字"
        sleep 3s
        start_menu
        ;;
    esac
}

start_menu
