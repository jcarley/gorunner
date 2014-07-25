
-- +goose Up
--  ALTER TABLE `gorunner`.`job_tasks` ADD INDEX `job_id` USING BTREE (`job_id`) comment '';
--  ALTER TABLE `gorunner`.`job_tasks` ADD INDEX `task_id` USING BTREE (`task_id`) comment '';

ALTER TABLE `gorunner`.`jobs` ADD CONSTRAINT `fk_jobs_job_tasks_1` FOREIGN KEY (`id`) REFERENCES `gorunner`.`job_tasks` (`job_id`)  ON DELETE CASCADE;
--  ALTER TABLE `gorunner`.`job_tasks` ADD CONSTRAINT `fk_jobs_job_tasks_2` FOREIGN KEY (`task_id`) REFERENCES `gorunner`.`tasks` (`id`) ON DELETE CASCADE;


-- +goose Down
ALTER TABLE `gorunner`.`jobs` DROP FOREIGN KEY `fk_jobs_job_tasks_1`;
--  ALTER TABLE `gorunner`.`job_tasks` DROP INDEX `job_id`, DROP INDEX `task_id`;

