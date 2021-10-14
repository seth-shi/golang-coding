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

```
pipeline {
  agent any
  stages {
    stage('检出') {
      steps {
        checkout([$class: 'GitSCM',
        branches: [[name: GIT_BUILD_REF]],
        userRemoteConfigs: [[
          url: GIT_REPO_URL,
          credentialsId: CREDENTIALS_ID
        ]]])
      }
    }
    stage('设置环境变量') {
      steps {
        writeFile(file: '.env', text: 'DATABASE_DSN=mysql:root:12345678@tcp(127.0.0.1:3306)/test2')
      }
    }
    stage('构建镜像并推送到 CODING Docker 制品库') {
      steps {
        sh "docker build -t ${CODING_DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION} -f ${DOCKERFILE_PATH} ${DOCKER_BUILD_CONTEXT}"
        useCustomStepPlugin(key: 'codingcorp:artifact_docker_push', version: 'latest', params: [image:"${CODING_DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}",repo:"${DOCKER_REPO_NAME}"])
      }
    }
    stage('部署到远端服务') {
      steps {
        script {
          def remoteConfig = [:]
          remoteConfig.name = "my-remote-server"
          remoteConfig.host = "${REMOTE_HOST}"
          remoteConfig.port = "${REMOTE_SSH_PORT}".toInteger()
          remoteConfig.allowAnyHosts = true

          withCredentials([
            sshUserPrivateKey(
              credentialsId: "${REMOTE_CRED}",
              keyFileVariable: "privateKeyFilePath"
            ),
            usernamePassword(
              credentialsId: "${CODING_ARTIFACTS_CREDENTIALS_ID}",
              usernameVariable: 'CODING_DOCKER_REG_USERNAME',
              passwordVariable: 'CODING_DOCKER_REG_PASSWORD'
            )
          ]) {
            // SSH 登陆用户名
            remoteConfig.user = "${REMOTE_USER_NAME}"
            // SSH 私钥文件地址
            remoteConfig.identityFile = privateKeyFilePath

            // 请确保远端环境中有 Docker 环境
            sshCommand(
              remote: remoteConfig,
              command: "docker login -u ${CODING_DOCKER_REG_USERNAME} -p ${CODING_DOCKER_REG_PASSWORD} ${CODING_DOCKER_REG_HOST}",
              sudo: true,
            )

            sshCommand(
              remote: remoteConfig,
              command: "docker rm -f golang-coding | true",
              sudo: true,
            )

            // DOCKER_IMAGE_VERSION 中涉及到 GIT_LOCAL_BRANCH / GIT_TAG / GIT_COMMIT 的环境变量的使用
            // 需要在本地完成拼接后，再传入到远端服务器中使用
            DOCKER_IMAGE_URL = sh(
              script: "echo ${CODING_DOCKER_REG_HOST}/${CODING_DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}",
              returnStdout: true
            )

            sshCommand(
              remote: remoteConfig,
              command: "docker run -d -p 8199:8199 --name golang-coding ${DOCKER_IMAGE_URL}",
              sudo: true,
            )

            // 删除没用使用的镜像
            sshCommand(
              remote: remoteConfig,
              command: "docker image prune -a -f",
              sudo: true,
            )

            echo "部署成功，请到 http://${REMOTE_HOST}:8199 预览效果"
          }
        }

      }
    }
  }
  environment {
    CODING_DOCKER_REG_HOST = "${CCI_CURRENT_TEAM}-docker.pkg.${CCI_CURRENT_DOMAIN}"
    CODING_DOCKER_IMAGE_NAME = "${PROJECT_NAME.toLowerCase()}/${DOCKER_REPO_NAME}/${DOCKER_IMAGE_NAME}"
  }
}
```

### 项目仓库参考
[https://github.com/seth-shi/golang-coding](https://github.com/seth-shi/golang-coding)