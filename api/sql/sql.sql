CREATE DATABASE IF NOT EXISTS socialmedia;
USE socialmedia;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id int auto_increment PRIMARY KEY,
    name varchar(50) NOT NULL,
    nick varchar(50) NOT NULL UNIQUE,
    email varchar(50) NOT NULL UNIQUE,
    pass varchar(100) NOT NULL,
    createdat TIMESTAMP default current_timestamp()
) ENGINE=INNODB;
