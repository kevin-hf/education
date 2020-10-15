### 环境依赖：

- **Ubuntu 16.04**
- **git**
- **docker 17.03.0-ce+**
- **docker-compose 1.8**
- **Golang 1.11.x+**
- **make**

### 安装步骤

1. 新建目录

   ```shell
   $ mkdir -p $GOPATH/src/github.com/kongyixueyuan.com/
   $ cd $GOPATH/src/github.com/kongyixueyuan.com/
   ```

2. 下载项目

   ```go
   $ git clone https://github.com/kevin-hf/education.git
   ```

3. 进入fixtures目录启动网络

   ```shell
   $ cd $GOPATH/src/github.com/kongyixueyuan.com/education/fixtures
   $ docker-compose up
   ```

4. 返回至项目根目录

   ```shell
   $ cd $GOPATH/src/github.com/kongyixueyuan.com/education
   ```

5. 编译

   ```shell
   $ go build
   ```

6. 启动服务

   ```shell
   $ ./education
   ```

7. 浏览器访问

   ```url
   http://localhost:9000
   ```

8. 登录用户名及密码

   ```
   用户名：Hanxiaodong
   密码：123456
   ```

9. 停止服务

   Ctrl + C 停止Web服务，然后使用如下命令清空：

   ```shell
   $ make clean
   ```

   ​