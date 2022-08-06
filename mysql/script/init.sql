CREATE TABLE attendances(
    id int auto_increment,
    user_id VARCHAR(60),
    start_time VARCHAR(40),
    end_time VARCHAR(40) DEFAULT NULL,
    working_time VARCHAR(40) DEFAULT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    primary key (id)
);

