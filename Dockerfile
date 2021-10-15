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
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go


###############################################################################
#          运行镜像: 从构建镜像中复制构建好的文件
###############################################################################
# 极客可以尝试 scratch
FROM scratch


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
# scratch 镜像需要使用 CMD ["./main"] 方式运行
###############################################################################
WORKDIR $WORKDIR
CMD ["./main"]