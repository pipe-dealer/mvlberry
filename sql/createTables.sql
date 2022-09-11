-- \i 'C:/Users/test/mvlberry/sql/createTables.sql'

-- DROP TABLE IF EXISTS users, friendships, requests CASCADE; 

-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     username TEXT UNIQUE NOT NULL,
--     password VARCHAR(32) NOT NULL /* MD5 hash*/
-- );

-- CREATE TABLE friendships (
--     id INTEGER NOT NULL references "users" (id),
--     f_id INTEGER NOT NULL references "users" (id),
--     fs_id INTEGER

-- );

-- CREATE TABLE requests (
--     req_id SERIAL PRIMARY KEY,
--     id INTEGER NOT NULL references "users" (id),
--     r_id INTEGER NOT NULL references "users" (id)
-- );

CREATE TABLE messages (
    msg_id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL references "users" (id),
    receiver_id INTEGER NOT NULL references "users" (id),
    msg_text TEXT NOT NULL,
    date_sent TIMESTAMP NOT NULL
);