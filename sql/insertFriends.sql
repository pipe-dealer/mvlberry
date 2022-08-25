-- \i 'C:/Users/test/mvlberry/sql/insertFriends.sql'

-- test1 - test2,test3,test5
-- test2 - test1,test3,test4
-- test3 - test1,test2,test4
-- test4 - test2,test3,test5
-- test5 - test1,test4
INSERT INTO friendships (id,f_id,fs_id)
VALUES
(1,2,1),(2,1,1),

(1,3,2),(3,1,2),

(1,5,3),(5,1,3),

(2,3,4),(3,2,4),

(2,4,5),(4,2,5),

(3,4,6),(4,3,6),

(4,5,7),(5,4,7)