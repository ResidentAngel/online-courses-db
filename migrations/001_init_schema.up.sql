-- Справочники и независимые таблицы
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       login VARCHAR(50) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       role VARCHAR(20) NOT NULL CHECK (role IN ('student', 'teacher', 'company', 'admin'))
);

CREATE TABLE userentry (
                           id SERIAL PRIMARY KEY,
                           user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           entry_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           action VARCHAR(100) NOT NULL
);

CREATE TABLE language (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(50) NOT NULL,
                          compiler_version VARCHAR(50) NOT NULL
);

CREATE TABLE company (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(100) UNIQUE NOT NULL,
                         description TEXT
);

CREATE TABLE page (
                      id SERIAL PRIMARY KEY,
                      title VARCHAR(200) NOT NULL,
                      content TEXT NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE problem (
                         id SERIAL PRIMARY KEY,
                         title VARCHAR(200) NOT NULL,
                         description TEXT NOT NULL,
                         difficulty VARCHAR(20) CHECK (difficulty IN ('easy', 'medium', 'hard')),
                         time_limit_ms INT NOT NULL DEFAULT 1000 CHECK (time_limit_ms <= 5000),
                         memory_limit_mb INT NOT NULL DEFAULT 64
);

CREATE TABLE tips (
                      id SERIAL PRIMARY KEY,
                      problem_id INT NOT NULL REFERENCES problem(id) ON DELETE CASCADE,
                      hint_text TEXT NOT NULL
);

-- Связующие таблицы M:N
CREATE TABLE languagetoproblem (
                                   problem_id INT NOT NULL REFERENCES problem(id) ON DELETE CASCADE,
                                   language_id INT NOT NULL REFERENCES language(id) ON DELETE CASCADE,
                                   PRIMARY KEY (problem_id, language_id)
);

CREATE TABLE problem_to_company (
                                    problem_id INT NOT NULL REFERENCES problem(id) ON DELETE CASCADE,
                                    company_id INT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
                                    PRIMARY KEY (problem_id, company_id)
);

-- Проверка кода (с проверкой бизнес-правила времени выполнения <= 5000 мс)
CREATE TABLE coderun (
                         id SERIAL PRIMARY KEY,
                         user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         problem_id INT NOT NULL,
                         language_id INT NOT NULL,
                         source_code TEXT NOT NULL,
                         output TEXT,
                         exec_time_ms INT CHECK (exec_time_ms <= 5000),
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (problem_id, language_id) REFERENCES languagetoproblem(problem_id, language_id)
);

CREATE TABLE codesubmit (
                            id SERIAL PRIMARY KEY,
                            user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            problem_id INT NOT NULL,
                            language_id INT NOT NULL,
                            source_code TEXT NOT NULL,
                            status VARCHAR(30) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'in_progress', 'accepted', 'wrong_answer', 'time_limit')),
                            exec_time_ms INT CHECK (exec_time_ms <= 5000),
                            memory_used_mb INT,
                            submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Составной внешний ключ гарантирует бизнес-правило 7.1:
    -- отправить код можно только на языке, разрешенном для этой задачи!
                            FOREIGN KEY (problem_id, language_id) REFERENCES languagetoproblem(problem_id, language_id)
);

-- Тестирование
CREATE TABLE test (
                      id SERIAL PRIMARY KEY,
                      title VARCHAR(200) NOT NULL,
                      description TEXT
);

CREATE TABLE testquestion (
                              id SERIAL PRIMARY KEY,
                              test_id INT NOT NULL REFERENCES test(id) ON DELETE CASCADE,
                              question_text TEXT NOT NULL
);

CREATE TABLE testanswer (
                            id SERIAL PRIMARY KEY,
                            question_id INT NOT NULL REFERENCES testquestion(id) ON DELETE CASCADE,
                            answer_text TEXT NOT NULL,
                            is_correct BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE teststart (
                           id SERIAL PRIMARY KEY,
                           user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           test_id INT NOT NULL REFERENCES test(id) ON DELETE CASCADE,
                           started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           UNIQUE (user_id, test_id)
);

CREATE TABLE testresult (
                            id SERIAL PRIMARY KEY,
                            user_id INT NOT NULL,
                            test_id INT NOT NULL,
                            score INT NOT NULL CHECK (score >= 0),
                            duration_sec INT NOT NULL,
                            completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Составной внешний ключ гарантирует бизнес-правило 7.1:
    -- результат теста невозможен без записи в teststart для этого пользователя и теста!
                            FOREIGN KEY (user_id, test_id) REFERENCES teststart(user_id, test_id) ON DELETE CASCADE
);

-- Биллинг
CREATE TABLE transactiontype (
                                 id SERIAL PRIMARY KEY,
                                 name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE "transaction" (
                               id SERIAL PRIMARY KEY,
                               user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               type_id INT NOT NULL REFERENCES transactiontype(id) ON DELETE RESTRICT,
                               amount NUMERIC(10, 2) NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);