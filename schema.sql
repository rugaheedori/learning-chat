CREATE DATABASE IF NOT EXISTS chatting;
USE chatting;

-- room 관련 테이블
CREATE TABLE room (
    id bigint primary key NOT NULL auto_increment,
    name varchar(255) NOT NULL UNIQUE,
    createdAt timestamp DEFAULT CURRENT_TIMESTAMP,
    updatedAt timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- chat 관련 테이블
CREATE TABLE chat (
    id bigint primary key NOT NULL auto_increment,
    room varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    message varchar(255) NOT NULL,
    `when` timestamp DEFAULT CURRENT_TIMESTAMP
);

-- server 관리 테이블
CREATE TABLE serverInfo (
    ip varchar(255) primary key NOT NULL,
    available bool NOT NULL
);
