-- CREATE game_info table
CREATE TABLE game_info (
    id INTEGER PRIMARY KEY,
    name VARCHAR(15) NOT NULL,
    data BLOB NOT NULL,
    valid_date DATE NOT NULL  
);