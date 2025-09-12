## 启动项目


### docker
```bash
 ENV=testShop docker compose --env-file .env.testshop up --build
```

### 本地
```bash
SHOP_ENV=development go run .
```


## 其他补充

```bash
docker ps # 查看容器列表 如 membership-backend

docker exec -it membership-backend /bin/sh # 进入容器
```