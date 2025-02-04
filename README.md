Тестовое задание:

Дана таблица categories, содержащая древовидную структуру каталога интернет-магазина. У каждой категории есть уникальный идентификатор (id) и ссылка на родительскую категорию (parent_id). Если parent_id равно NULL, то категория является корневой.

Пример:
-- Создание таблицы
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INT REFERENCES categories(id) ON DELETE CASCADE
);

-- Вставка данных
INSERT INTO categories (id, name, parent_id) VALUES
    (1, 'Электроника', NULL),
    (2, 'Телефоны', 1),
    (3, 'Ноутбуки', 1),
    (4, 'Смартфоны', 2),
    (5, 'Аксессуары', 2),
    (6, 'Чехлы', 5),
    (7, 'Зарядки', 5),
    (8, 'Одежда', NULL),
    (9, 'Мужская', 8),
    (10, 'Женская', 8);

-- Установка значения для последовательности
SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));


Необходимо реализовать эндпоинт, который будет получать на вход параметр с ID категории, возвращает все её дочерние элементы (включая вложенные уровни). Например, при вводе id = 2 (категория “Телефоны”) запрос должен вернуть:

id name  parent_id
4 Смартфоны 2
5 Аксессуары 2
6 Чехлы  5
7 Зарядки  5

Также необходимо получать потомков одним запросом к БД, можно модифицировать для этой задачи исходную таблицу, использовать триггеры при необходимости.
