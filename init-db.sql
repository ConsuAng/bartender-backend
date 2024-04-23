CREATE DATABASE IF NOT EXISTS angelusbartender;

USE angelusbartender;

CREATE TABLE IF NOT EXISTS user (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_cocktail (
    user_cocktail_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    cocktail_id VARCHAR(255) NOT NULL,
    cocktail_name VARCHAR(255) NOT NULL,
    cocktail_image VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id) 
        REFERENCES user(user_id)
        ON DELETE CASCADE
);