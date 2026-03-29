-- +goose Up
-- +goose StatementBegin

INSERT INTO staff.employees (id, full_name, position, created_by_id, created_by_name)
VALUES (1, 'Администратор', 1, 1, 'seed_migration')
ON CONFLICT (id) DO NOTHING;

INSERT INTO staff.offices (name, address, created_by_id, created_by_name)
VALUES
    ('Офис Москва Центр', 'г. Москва, ул. Тверская, д. 10', 1, 'seed_migration'),
    ('Офис Санкт-Петербург Невский', 'г. Санкт-Петербург, Невский проспект, д. 28', 1, 'seed_migration'),
    ('Офис Новосибирск Север', 'г. Новосибирск, Красный проспект, д. 101', 1, 'seed_migration'),
    ('Офис Екатеринбург Парк', 'г. Екатеринбург, ул. Малышева, д. 51', 1, 'seed_migration'),
    ('Офис Казань Кремль', 'г. Казань, ул. Баумана, д. 15', 1, 'seed_migration'),
    ('Офис Нижний Новгород Волга', 'г. Нижний Новгород, ул. Большая Покровская, д. 32', 1, 'seed_migration'),
    ('Офис Челябинск Восток', 'г. Челябинск, проспект Ленина, д. 64', 1, 'seed_migration'),
    ('Офис Самара Центральный', 'г. Самара, Московское шоссе, д. 17', 1, 'seed_migration'),
    ('Офис Ростов-на-Дону Юг', 'г. Ростов-на-Дону, ул. Большая Садовая, д. 82', 1, 'seed_migration'),
    ('Офис Уфа Башкортостан', 'г. Уфа, проспект Октября, д. 4', 1, 'seed_migration');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM staff.offices
WHERE created_by_name = 'seed_migration'
  AND name IN (
      'Офис Москва Центр',
      'Офис Санкт-Петербург Невский',
      'Офис Новосибирск Север',
      'Офис Екатеринбург Парк',
      'Офис Казань Кремль',
      'Офис Нижний Новгород Волга',
      'Офис Челябинск Восток',
      'Офис Самара Центральный',
      'Офис Ростов-на-Дону Юг',
      'Офис Уфа Башкортостан'
  );

DELETE FROM staff.employees
WHERE id = 1
  AND full_name = 'Администратор'
  AND created_by_name = 'seed_migration';

-- +goose StatementEnd
