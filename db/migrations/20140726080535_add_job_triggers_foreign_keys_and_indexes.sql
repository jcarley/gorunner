
-- +goose Up
ALTER TABLE `gorunner`.`job_triggers` ADD INDEX `job_id` USING BTREE (`job_id`) comment '';
ALTER TABLE `gorunner`.`job_triggers` ADD INDEX `trigger_id` USING BTREE (`trigger_id`) comment '';

ALTER TABLE `gorunner`.`job_triggers` ADD CONSTRAINT `fk_jobs_job_triggers_1` FOREIGN KEY (`job_id`) REFERENCES `gorunner`.`jobs` (`id`)  ON DELETE CASCADE;
ALTER TABLE `gorunner`.`job_triggers` ADD CONSTRAINT `fk_triggers_job_triggers_1` FOREIGN KEY (`trigger_id`) REFERENCES `gorunner`.`triggers` (`id`) ON DELETE CASCADE;


-- +goose Down
ALTER TABLE `gorunner`.`job_triggers` DROP FOREIGN KEY `fk_jobs_job_triggers_1`, DROP FOREIGN KEY `fk_triggers_job_triggers_1`;
ALTER TABLE `gorunner`.`job_triggers` DROP INDEX `job_id`, DROP INDEX `trigger_id`;

