CREATE TABLE `payment_status` (
  `payment_id` int NOT NULL,
  `status_id` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER TABLE `payment_status` ADD FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`);
ALTER TABLE `payment_status` ADD FOREIGN KEY (`status_id`) REFERENCES `statuses` (`id`);