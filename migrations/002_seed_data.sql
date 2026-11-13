-- Начальные данные для демонстрации
INSERT INTO users (login, password_hash, email, role) VALUES
                                                          ('student_ivan', '$2a$10$hash1234567890', 'ivan@uni.ru', 'student'),
                                                          ('teacher_petr', '$2a$10$hash0987654321', 'petr@uni.ru', 'teacher'),
                                                          ('admin_root', '$2a$10$hashadminroot1', 'admin@uni.ru', 'admin');

INSERT INTO language (name, compiler_version) VALUES
                                                  ('Go', '1.22'),
                                                  ('Python', '3.12'),
                                                  ('C++', 'GCC 13.2');

INSERT INTO company (name, description) VALUES
                                            ('Яндекс', 'Задачи с алгоритмических собеседований'),
                                            ('Т-Банк', 'Финтех-задачи и структуры данных');

INSERT INTO page (title, content) VALUES
    ('Введение в алгоритмы и О-большое', 'Сложность алгоритма оценивается по времени и памяти...');

INSERT INTO problem (title, description, difficulty, time_limit_ms, memory_limit_mb) VALUES
                                                                                         ('Сумма двух чисел (Two Sum)', 'Найдите индексы двух чисел в массиве, сумма которых равна target.', 'easy', 1000, 64),
                                                                                         ('Сбалансированные скобки', 'Проверьте корректность скобочной последовательности.', 'medium', 1500, 32);

INSERT INTO tips (problem_id, hint_text) VALUES
                                             (1, 'Используйте хеш-таблицу (map) для поиска за O(N).'),
                                             (2, 'Используйте стек (stack) для хранения открывающих скобок.');

INSERT INTO languagetoproblem (problem_id, language_id) VALUES (1, 1), (1, 2), (2, 1), (2, 3);
INSERT INTO problem_to_company (problem_id, company_id) VALUES (1, 1), (2, 2);

INSERT INTO test (title, description) VALUES ('Базовый тест по SQL и БД', 'Проверка знаний нормальных форм и транзакций');
INSERT INTO testquestion (test_id, question_text) VALUES (1, 'Какая команда завершает транзакцию с сохранением изменений?');
INSERT INTO testanswer (question_id, answer_text, is_correct) VALUES (1, 'COMMIT', true), (1, 'ROLLBACK', false);

INSERT INTO transactiontype (name) VALUES ('Пополнение баланса'), ('Покупка PRO-доступа'), ('Открытие платной подсказки');
INSERT INTO "transaction" (user_id, type_id, amount) VALUES (1, 1, 1500.00), (1, 2, -990.00);