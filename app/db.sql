create database tickets;
use tickets;
CREATE TABLE tickets (
id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
t_title VARCHAR(100) NOT NULL,
t_desc VARCHAR(5000) NOT NULL,
t_status VARCHAR(10),
t_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
INSERT INTO tickets (t_title, t_desc, t_status) VALUES ("Cannot access inventory app", "When openeing the application if gives me access denied even when I am using the correct credentials", "open");
INSERT INTO tickets (t_title, t_desc, t_status) VALUES ("Error while saving the purchases", "When trying to save the purchased item, the status is not updated", "open");
INSERT INTO tickets (t_title, t_desc, t_status) VALUES ("Need access to the customer relation app","Justification: need this access because of job requirements", "open");
INSERT INTO tickets (t_title, t_desc, t_status) VALUES ("Error while saving the requests", "Error message: 503 when trying to save", "open");
INSERT INTO tickets (t_title, t_desc, t_status) VALUES ("Cannot retreive customer data","The result is an empty page", "open");