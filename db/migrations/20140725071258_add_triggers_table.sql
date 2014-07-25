
-- +goose Up
SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `triggers` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(125) DEFAULT NULL,
  `schedule` varchar(255) DEFAULT NULL,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  `version` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `triggers`;
SET FOREIGN_KEY_CHECKS = 1;
-- SQL section 'Down' is executed when this migration is rolled back

