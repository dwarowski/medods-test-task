Запуск API без docker
```bash
git clone git@gitgub.com:dwarowski/medods-test-task

cd medods-test-task

source .env

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
POSTGRES_HOST="localhost"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD=12345678
POSTGRES_DB="postgres"
POSTGRES_PORT=5432

PRIVATE_KEY_PATH="keys/private.pem"
PUBLIC_KEY_PATH="keys/public.pem"

GIN_MODE="debug"
```