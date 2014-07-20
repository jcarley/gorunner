
-- +goose Up
SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `job_tasks` (
  `job_id` bigint(20) NOT NULL,
  `task_id` bigint(20) NOT NULL,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `job_tasks`;
SET FOREIGN_KEY_CHECKS = 1;
-- SQL section 'Down' is executed when this migration is rolled back

