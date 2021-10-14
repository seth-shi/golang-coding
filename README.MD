# Golang Ci

* [https://coding.net/](https://coding.net/)每个月有`1000`分钟的免费构建资源
使用之前先查看一遍文档: [https://help.coding.net/docs/ci/start.html](https://help.coding.net/docs/ci/start.html)

### Start
* 使用的[https://github.com/gogf/gf](https://github.com/gogf/gf)开发, 参考它的默认`Dockerfile`并修改了部分(默认的需要在宿主机打包[为了依赖下载缓存], 我修改为全部在镜像中进行)

```
### Dockerfile
###############################################################################
#          使用第一个镜像打包
###############################################################################
FROM golang as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 设置固定的打包路径
WORKDIR /app
COPY . .

# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod tidy && go build -o main main.go


###############################################################################
#          运行镜像: 从构建镜像中复制构建好的文件
###############################################################################
# 极客可以尝试 scratch
# FROM scratch
FROM alpine

# 使用国内alpine源
RUN echo http://mirrors.aliyun.com/alpine/v3.8/main/ > /etc/apk/repositories

# 设置系统时区 - +8时区
RUN apk update && apk add tzdata ca-certificates bash
RUN rm -rf /etc/localtime && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone


# 设置固定的项目路径
ENV WORKDIR /var/www

# 添加应用可执行文件，并设置执行权限
COPY --from=builder /app/main $WORKDIR/main

# 添加I18N多语言文件、静态文件、配置文件、模板文件
COPY --from=builder /app/i18n     $WORKDIR/i18n
COPY --from=builder /app/public   $WORKDIR/public
COPY --from=builder /app/config   $WORKDIR/config
COPY --from=builder /app/template $WORKDIR/template
COPY --from=builder /app/.env     $WORKDIR/.env

###############################################################################
#                                   运行
###############################################################################
WORKDIR $WORKDIR
CMD ["./main"]
```

* `Coding`新建构建计划的时候, 可以选择现有的模板,这里我选择`CODING Docker 镜像推送`
* 简单的如下使用
