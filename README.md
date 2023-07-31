# goads - URL shortener with ads

[RUS](#цель)

## Goal

Studying the process of building an application on a microservice architecture on Go using Docker

## Main idea

User can make his link shorter and add his ads to it.
The ad will be randomly selected from added ads and then shown to visitors before redirecting

## Setup

`docker-compose run`

## Структура проекта

1. Auth - user authentication. The database stores user data,
   which are used for password authentication, after which a JWT token is issued.
   The token is signed by the RSA algorithm, and is used for further authentication
2. Ads - work with ads. CRUD with the ability to filter ads by parameters:
    - Published (true/false)
    - Title (search)
    - Author (id)
    - Date of creation
3. URL Shortener - shortening links. CRUD associated with Ads to receive ads when requesting a redirect.
   From the published ads, one is selected randomly
4. API Gateway is an intermediate stage between the user and microservices.
   Before some requests, performs authentication by accessing Auth

## Used technologies

- Docker, docker-compose
- PostgreSQL
- gRPC API
- REST API
- JWT
- Gin

---

## Цель

Изучение процесса построения приложения на микросервисной архитектуре на Go с использованием Docker

## Основная идея

Пользователь может сделать свою ссылку короче и добавить рекламу по ней
Реклама будет выбрана случайно из добавленных, а затем показана пользователю перед редиректом

## Запуск

`docker-compose up`

## Структура проекта

1. Auth - аутентификация пользователей. В БД хранятся пользовательские данные,
   по которым производится аутентификация по паролю, после чего выдается JWT токен.
   Токен подписывается алгоритмом RSA, и используется для дальнейшей аутентификации
2. Ads - работа с объявлениями. CRUD с возможностью фильтрации объявлений по параметрам:
    - Опубликовано (true/false)
    - Название (поиск)
    - Автор (id)
    - Дата создания
3. URL Shortener - сокращение ссылок. CRUD, связанный с Ads для получения объявлений при запросе редиректа.
   Из опубликованных объявлений выбирается одно случайным образом
4. API Gateway - промежуточный этап между пользователем и микросервисами.
   Перед некоторыми запросами совершает аутентификацию, обращаясь к Auth

## Использованные технологии

- Docker, docker-compose
- PostgreSQL
- gRPC API
- REST API
- JWT
- Gin
