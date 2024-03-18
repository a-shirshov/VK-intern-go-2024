## Intern Задание

### Запуск
    docker compose up --build       

## Сваггер документация
Документация доступна по /docs

## Команды
Посмотреть покрытие - make cover. Покрытие - 70%

Запустить автоматическую генерацию документации - make swagger

## Авторизация
Сделано очень просто. При сборке проекте в базе под ID 2 - администратор.

Для того чтобы воспользоваться методами добавления/изменения:

В header Authorization поставить id администратора - в нашем случае 2

В swagger - есть кнопка для авторизации. 

Туда также можно написать 2, чтобы воспользоваться методами

## Еще моменты
Проект сделан по чистой архитектуре

Даты в формате - YYYY-MM-DD

Для тестирования использовался - gomock, pgxmock

Для работы с БД - pgx

Для документации - go-swagger

Порт по умолчанию - 8080