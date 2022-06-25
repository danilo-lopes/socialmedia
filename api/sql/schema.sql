CREATE DATABASE IF NOT EXISTS socialmedia;
USE socialmedia;

DROP TABLE IF EXISTS publications;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    nick VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    pass VARCHAR(100) NOT NULL,
    createdat TIMESTAMP DEFAULT current_timestamp()
) ENGINE=INNODB;

CREATE TABLE followers(
    user_id INT NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    follower_id INT NOT NULL,
        FOREIGN KEY(follower_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, follower_id)
) ENGINE=INNODB;

CREATE TABLE publications(
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    content VARCHAR(300) NOT NULL,
    author_id INT NOT NULL,
        FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE CASCADE,
    likes INT DEFAULT 0,
    createdat TIMESTAMP DEFAULT current_timestamp()
) ENGINE=INNODB;