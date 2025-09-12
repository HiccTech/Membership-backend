## 启动项目


### docker
```bash
 ENV=testShop docker compose --env-file .env.testshop up --build
```

### 本地
```bash
SHOP_ENV=development go run .
```