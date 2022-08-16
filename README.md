### 所需库
```powershell
go get -u gopkg.in/ini.v1
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/gin-gonic/gin
go get -u github.com/gin-contrib/sessions
go get -u github.com/gin-contrib/sessions/cookie
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/gin-gonic/gin
go get -u github.com/gin-contrib/cors
```

### TODO List
- [ ] 页面访问量排名
- [ ] PV UV数据请求接口
- [x] JsError数据请求接口
- [x] ApiError数据请求接口
- [ ] SourceError数据请求接口
- [ ] Performance数据请求接口
- [x] 资源GET请求参数解析中间件
- [ ] 使用channel将上报逻辑改为异步批量插入数据库
- [x] 生成测试数据脚本
- [x] PV UV数据上报
- [x] JWT中间件
- [x] 项目请求接口

[在线的时间戳转换工具](https://tool.lu/timestamp/)

### 创建数据库
```sql
CREATE DATABASE simple_monitor CHARACTER SET utf8mb4;
```

### 运行
```bash
nohup go run main.go > myout.file 2>&1 &
```