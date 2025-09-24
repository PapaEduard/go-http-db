Мини-проект Go project
Дано приложение на golang https://github.com/AnastasiyaGapochkina01/go-http-db, которое реализует простой HTTP-сервер с записью и просмотром истории запросов в PostgreSQL.

Необходимо организовать для него devops-конвейер:

Запустить в docker:
делать multistage сборку
публиковать собранный image в docker hub
у БД должен быть healthcheck
деплой осуществлять с помощью docker compose, который должен содержать сервисы
приложение
nginx (в роли reverse proxy)
postgres
migrator
migrator:
  image: postgres:15
  restart: on-failure
  environment:
    PGPASSWORD: postgres
  volumes:
    - ./migrations:/migrations
  command: >
    sh -c "while ! pg_isready -h db -U postgres; do sleep 1; done && psql -h db -U postgres -d app_db -f /migrations/001_init.sql"
CI/CD:
pipeline должен содержать этапы
linter (для файлов приложения, файлов docker)
build
сканирование на уязвимости image
test (запуск всей инфраструктуры любым способом и простой запрос curl http://$host:8080/hello)
deploy
Написать bash-скрипт для резервного копирования базы данных
