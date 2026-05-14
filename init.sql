CREATE TABLE IF NOT EXISTS student (
    stdid int NOT NULL,
    firstname varchar(45) NOT NULL,
    lastname varchar(45) DEFAULT NULL,
    email varchar(45) NOT NULL,
    PRIMARY KEY (stdid),
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS course (
    cid varchar(32) NOT NULL,
    coursename varchar(45) NOT NULL,
    PRIMARY KEY (cid)
);

CREATE TABLE IF NOT EXISTS enroll (
    std_id int NOT NULL,
    course_id varchar(45) NOT NULL,
    date_enrolled varchar(45) DEFAULT NULL,
    PRIMARY KEY (std_id, course_id),
    CONSTRAINT course_fk FOREIGN KEY (course_id) REFERENCES course (cid) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT std_fk FOREIGN KEY (std_id) REFERENCES student (stdid) ON DELETE CASCADE ON UPDATE CASCADE
);
