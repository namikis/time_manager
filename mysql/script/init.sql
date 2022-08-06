CREATE TABLE attendances(
    id int auto_increment,
    user_id VARCHAR(60),
    start_time VARCHAR(40),
    end_time VARCHAR(40) DEFAULT NULL,
    working_time VARCHAR(40) DEFAULT NULL,
    breaking_time VARCHAR(40) DEFAULT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    primary key (id)
);

CREATE TABLE breaks(
    id int auto_increment,
    attendance_id int,
    start_break_time VARCHAR(40),
    end_break_time VARCHAR(40) DEFAULT NULL,
    breaking_time VARCHAR(40) DEFAULT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    primary key (id),
    FOREIGN KEY fk_attendance REFERENCES attendances(id)
);

