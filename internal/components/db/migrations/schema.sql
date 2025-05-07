CREATE TABLE IF NOT EXISTS users (
                                     chat_id BIGINT PRIMARY KEY,
                                     start_date DATETIME
);

CREATE TABLE IF NOT EXISTS lessons (
                                       id INT AUTO_INCREMENT PRIMARY KEY,
                                       day_number INT NOT NULL,
                                       content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_lessons (
                                            id INT AUTO_INCREMENT PRIMARY KEY,
                                            user_id BIGINT,
                                            lesson_id INT,
                                            sent_at DATETIME,
                                            UNIQUE KEY unique_user_lesson (user_id, lesson_id)
    );