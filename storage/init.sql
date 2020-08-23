CREATE TABLE IF NOT EXISTS news (
  id  SERIAL  UNIQUE  NOT NULL PRIMARY KEY,
  title  text UNIQUE ,
  description text,
  link text
);
