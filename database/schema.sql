CREATE TABLE `orders` (
  `id` varchar(36) NOT NULL,
  `origin_lat` decimal(8,6) NOT NULL,
  `origin_long` decimal(9,6) NOT NULL,
  `destination_lat` decimal(8,6) NOT NULL,
  `destination_long` decimal(9,6) NOT NULL,
  `distance` int(11) NOT NULL,
  `status` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_list_orders` (`id`,`distance`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;