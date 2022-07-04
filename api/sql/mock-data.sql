/* 
"User1" "i7e8Iy9N5Y0yJ5VoK9Ej6bIA5Dgilm", "yaJSX2awOboJyxaT39CmIh1XLsA9YM"
"User2" "5A6CFCGYIheDUFI56fe6kTHF9AZqCP"
"User3" "ZnbcCM5zLjccxW4E9Oy1Jvo0j8RDME"
"User4" "8WPOaKok2C2UNUX4HxTEEpPm7XTI82"
*/

INSERT INTO users(name, nick, email, pass) VALUES
    ("User1", "usr1", "user1@gmail.com", "$2a$10$BRIzh6AL/NoSnXSbhMxh7OMBrwcPDi5fOm6xz6uC6snirTYQ3Vt/e"),
    ("User2", "usr2", "user2@gmail.com", "$2a$10$qPBFuNQmO0/RXieL.KXyA.fBOrlreu1k2kvaxcKvFR9OZ4OxUlwIm"),
    ("User3", "usr3", "user3@gmail.com", "$2a$10$G2981IbtaH1o2S6Qi.PuYO4J3/ftPkbvOvRo66rNY4vrm/BGN9Hqe"),
    ("User4", "usr4", "user4@gmail.com", "$2a$10$Vyi2hhqIP3vVBQBJarJjVe60Lt3paXLetBB2W2zCofyX3c5SzM7Ve");

INSERT INTO followers(user_id, follower_id) VALUES
    (1, 2),
    (1, 3),
    (3, 2),
    (4, 1),
    (2, 1),
    (3, 4);

INSERT INTO publications(title, content, author_id) VALUES
    ("User1 Publication", "This is the publication of User 1 !, oooohh", 1),
    ("User2 Publication", "This is the publication of User 2 !, oooohh", 2),
    ("User3 Publication", "This is the publication of User 3 !, oooohh", 3);
