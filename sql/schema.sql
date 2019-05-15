CREATE TABLE IF NOT EXISTS facts (
  user_id VARCHAR(46) PRIMARY KEY,
  facts TEXT
);

CREATE TABLE IF NOT EXISTS profiles (
  user_id VARCHAR(46) PRIMARY KEY,
  handle VARCHAR (50) UNIQUE,
  first_name VARCHAR(50),
  last_name VARCHAR (50),
  email VARCHAR (255),
  phone VARCHAR (15),
  pic_url VARCHAR (255)
);
