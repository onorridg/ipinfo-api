# IP INFO API - cервис для получения информации об IP адресе
![ip-info-api](https://user-images.githubusercontent.com/83474704/208415009-d3299b27-f4cb-434e-ae32-71d33c84cfbd.png)

## Обязательно измените поля .ENV
```bash
JWT_SECRET_KEY=
PASSWORD_SALT=        
SWAGGER_DOCS_HOST=   # Укажите свой хост для swagger docs
```
## Установка
```bash
git clone https://github.com/onorridg/ipinfo-api
cd ipinfo-api
make 
```
## Управление Docker
```bash
make stop   # Остановить api и redisDB
make start  # Запустить  api и redisDB
```
## Документация
**Авторизация: `Bearer Token`**

**https://ip.onorridg.tech/swagger/index.html**

![Screenshot swagger docs](https://user-images.githubusercontent.com/83474704/208075081-93840301-8162-46cf-b0f3-652c93df1e87.png)


## Пример запроса:
![Screenshot api request](https://user-images.githubusercontent.com/83474704/208075834-80007709-82ac-4956-907c-df2d2f631462.png)

## Протестировать API:
```bash
https://ip.onorridg.tech/api/v1   # Base URL
```

