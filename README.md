# ChatGPT api 服务
## 特性
1. 适配所有能用"自定义OpenAI域名"的客户端或网页
2. 不需要把OpenAI的API key放在客户端或网页上，使得API key不会被盗用
3. 支持用户管理功能，为每个用户分配key
4. 统计每个用户的API调用次数

## 部署
1. 创建配置文件config.yaml，配置以下内容：
```yaml
GinPort: 8080
OpenAIKey: "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
DBName: "chat.db"
```

2. 创建一个空的数据库文件chat.db:
```bash
touch chat.db
```

3. 运行docker:
```bash
docker run --name=chatapi -d \
  --restart=unless-stopped -p 8080:8080 \
  -v /root/chat/config.yaml:/web/config.yaml \
  -v /root/chat/chat.db:/web/chat.db \
  libli/chat:latest
```
把上面命令中的/root/chat替换为你的配置文件和数据库文件所在的目录。

4. 创建用户：
确保已经安装sqlite3客户端，然后运行：
```bash
sqlite3 chat.db
sqlite> insert into users (username, token) VALUES ('***', '****');
```

## 开源协议
MIT，随便拿去用，记得多帮我宣传宣传。

如果觉得帮助到你了，欢迎请[我喝一杯咖啡](https://github.com/libli/buy-me-a-coffee) ☕️。