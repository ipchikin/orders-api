CREATE TABLE orders (
    id int,
    `origin` POINT NOT NULL,
    `destination` POINT NOT NULL,
	distance int,
	status varchar(255),
    SPATIAL INDEX `SPATIAL` (`origin`)
) ENGINE=InnoDB;