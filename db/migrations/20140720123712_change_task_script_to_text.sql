
-- +goose Up
ALTER TABLE `gorunner`.`tasks` CHANGE COLUMN `script` `script` text CHARACTER SET utf8 DEFAULT NULL;
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
ALTER TABLE `gorunner`.`tasks` CHANGE COLUMN `script` `script` varchar(255) DEFAULT NULL;
-- SQL section 'Down' is executed when this migration is rolled back

