CREATE DATABASE IF NOT EXISTS `qframe`;
CREATE TABLE IF NOT EXISTS `qframe`.`mockdata`
(
    `int_col` integer NOT NULL, 
    `float_col` double precision,
    `string_col` text NOT NULL,
    `bool_col` boolean,
    PRIMARY KEY (`int_col`)
);
ALTER TABLE `qframe`.`mockdata` ADD INDEX (`string_col`(255));
GRANT ALL PRIVILEGES ON `qframe`.`mockdata` TO 'qframe';