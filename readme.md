# Требования
## Требования для запуска
* docker >= 28.3.0
* docker-compose >= 2.38.1

## Требования для разработки
* golang >= 1.24.4
* postgres >= 16.6

# Запуск через docker
## Шаг 1
Склонируйте текущий репозиторий себе на пк и перейдите в папку с ним
```bash
git clone git@gitgub.com:dwarowski/medods-test-task

cd medods-test-task
```

## Шаг 2
После чего откройте\создайте файл .env через любой удобный вам текстовый редактор и скопируйте туда следущие строчки изменив или оставив следующие поля:
* ``POSTGRES_USER``
* ``POSTGRES_PASSWORD``
* ``POSTGRES_DB``
* ``POSTGRES_PORT``

```bash
POSTGRES_HOST="db"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD=12345678
POSTGRES_DB="postgres"
POSTGRES_PORT=5432
GIN_MODE="release"
```

## Шаг 3
Если все предыдущие шаги получились успешно то запустите файл docker-compose следующей командой  
```bash
docker-compose -f docker-compose.yml up -d
```

# Запуск API без docker 
_С использованием переменных из примера .env (не забудьте поменять на ваши параметры базы данных)_
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

# .env example 
```bash
# Hostname is container name of database if using compose
export POSTGRES_HOST="localhost"
export POSTGRES_USER="postgres"
export POSTGRES_PASSWORD=12345678
export POSTGRES_DB="postgres"
export POSTGRES_PORT=5432

export PRIVATE_KEY_PATH="keys/private.pem"
export PUBLIC_KEY_PATH="keys/public.pem"

export GIN_MODE="debug"
```