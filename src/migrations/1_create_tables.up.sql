CREATE TABLE IF NOT EXISTS users (
  id VARCHAR (255) PRIMARY KEY UNIQUE NOT NULL,
  email VARCHAR (255) UNIQUE NOT NULL,
  username VARCHAR (255) UNIQUE NOT NULL,
  password VARCHAR (255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS bots (
  `id` VARCHAR (255) PRIMARY KEY UNIQUE NOT NULL,
  `user_id` VARCHAR (255) NOT NULL REFERENCES users(id),
  `token` VARCHAR (255) NOT NULL,
  `webhook` VARCHAR (255) NOT NULL DEFAULT '',
  `is_verified` BOOLEAN DEFAULT FALSE NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `devel` BOOLEAN NOT NULL DEFAULT false,
  `version` VARCHAR (10) NOT NULL DEFAULT '',
  `handlegroups` BOOLEAN DEFAULT TRUE NOT NULL,
  `handlebroadcast` BOOLEAN DEFAULT FALSE NOT NULL,
);
