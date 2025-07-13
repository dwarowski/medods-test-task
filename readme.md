Запуск API
```bash
git clone git@gitgub.com:dwarowski/medods-test-task

cd medods-test-task

go run main.go
```

Приложение и компоненты находятся в папке ``src/`` 
Запуск проиходит из корневой директории файла ``main.go`` 
командой 
```bash
go run main.go
```

.env example 
```bash
DATABASE_CONFIG="host=localhost user=postgres password=12345678 dbname=postgres port=5433 sslmode=disable"
PRIVATE_KEY_PATH="keys/private.pem"
PUBLIC_KEY_PATH="keys/public.pem"
GIN_MODE="debug"
```