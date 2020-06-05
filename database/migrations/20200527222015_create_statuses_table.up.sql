CREATE TABLE `statuses` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL
);

INSERT INTO `statuses` (`name`, `slug`) VALUES
("Ожидает оплаты", "expected"),
("Оплачено", "paid"),
("Ошибка", "error")