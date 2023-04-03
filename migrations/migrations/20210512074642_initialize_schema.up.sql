CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    identifier TEXT(255) NOT Null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_name TEXT(255),
    email TEXT(255) NOT NULL,
    CONSTRAINT unique_email unique (email(255))
);CREATE TABLE IF NOT EXISTS channels (
    id INT PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    title TEXT(255) NOT NULL,
    channel_name TEXT(255) NOT NULL,
    channel_secret TEXT(255),
    host_passphrase TEXT(255) NOT NULL,
    viewer_passphrase TEXT(255),
    recording_uid INT,
    recording_sid TEXT(255),
    recording_rid TEXT(255),
    dtmf TEXT(255)
);CREATE TABLE tokens (
    id INT PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    token_id TEXT(255),
    user_id INT,
    CONSTRAINT tokens_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);CREATE TABLE IF NOT EXISTS credentials (
    id INT PRIMARY KEY AUTO_INCREMENT,
    code TEXT(255) NOT NULL,
    access_token TEXT(255) NOT NULL,
    refresh_token TEXT(255) NOT NULL,
    token_type TEXT(255) NOT NULL,
    expiry TIMESTAMP
);