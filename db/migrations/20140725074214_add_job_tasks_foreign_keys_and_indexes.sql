
-- +goose Up
ALTER TABLE `gorunner`.`job_tasks` ADD INDEX `job_id` USING BTREE (`job_id`) comment '';
ALTER TABLE `gorunner`.`job_tasks` ADD INDEX `task_id` USING BTREE (`task_id`) comment '';

ALTER TABLE `gorunner`.`job_tasks` ADD CONSTRAINT `fk_jobs_job_tasks_1` FOREIGN KEY (`job_id`) REFERENCES `gorunner`.`jobs` (`id`)  ON DELETE CASCADE;
ALTER TABLE `gorunner`.`job_tasks` ADD CONSTRAINT `fk_tasks_job_tasks_1` FOREIGN KEY (`task_id`) REFERENCES `gorunner`.`tasks` (`id`) ON DELETE CASCADE;


-- +goose Down
ALTER TABLE `gorunner`.`job_tasks` DROP FOREIGN KEY `fk_jobs_job_tasks_1`, DROP FOREIGN KEY `fk_tasks_job_tasks_1`;
ALTER TABLE `gorunner`.`job_tasks` DROP INDEX `job_id`, DROP INDEX `task_id`;

