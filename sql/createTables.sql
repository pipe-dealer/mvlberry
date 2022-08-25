-- \i 'C:/Users/test/mvlberry/sql/createTables.sql'


DROP TABLE IF EXISTS users, friendships CASCADE; 

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password VARCHAR(32) NOT NULL /* MD5 hash*/
);

CREATE TABLE friendships (
    id INTEGER NOT NULL references "users" (id),
    f_id INTEGER NOT NULL references "users" (id),
    fs_id INTEGER

);