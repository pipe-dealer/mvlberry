DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS friendships;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password VARCHAR(32) NOT NULL /* MD5 hash*/
);

CREATE TABLE friendships (
    f_id SERIAL PRIMARY KEY,
    u1_id INTEGER NOT NULL references "users" (id),
    u2_id INTEGER NOT NULL references "users" (id)
);