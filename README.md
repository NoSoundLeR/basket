# Basket

## Идея
Один пользователь создаёт комнату, добавляет один вариант. Кидает ссылку на комнату другим. Другие заходят по ссылке, добавляют свой вариант. Создатель комнаты закрывает её. Из списка вариантов выбирается случайный и отсылается остальным по вебсокету или через poll. Пользователя проверяем по ip.

## Usage
```bash
go build -o ./basket-server ./cmd/server/server.go
./basket-server -bind 127.0.0.1:8080
```

## Example
![Simple](https://raw.githubusercontent.com/NoSoundLeR/basket/master/assets/simple.gif)

![Timeout](https://raw.githubusercontent.com/NoSoundLeR/basket/master/assets/timeout.gif)