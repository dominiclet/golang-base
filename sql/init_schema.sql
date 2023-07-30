CREATE DATABASE IF NOT EXISTS some_database;
USE some_database;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id` integer PRIMARY KEY AUTO_INCREMENT,
    `uuid` varchar(63) NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL,
    `password` varchar(255),
    `created_at` timestamp,
    `deleted_at` timestamp NULL,
    `updated_at` timestamp,
    `account_type` integer NOT NULL DEFAULT 0,
    `license_expiry` timestamp DEFAULT CURRENT_TIMESTAMP,
    `is_verified` int(1) NOT NULL DEFAULT 0,
    `verification_token` varchar(127)
);
CREATE UNIQUE INDEX user_uuid ON users (uuid);

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
    `id` integer PRIMARY KEY AUTO_INCREMENT,
    `user_id` integer,
    `token` varchar(255),
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `expires_at` timestamp
);

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
