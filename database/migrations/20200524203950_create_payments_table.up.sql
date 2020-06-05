CREATE TABLE `payments` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `uid` varchar(36),
  `amount` decimal(10, 2) NOT NULL,
  `purpose` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL
);