insert into role (name) values ('ROLE_ADMIN');
insert into role (name) values ('ROLE_COMPANY_OWNER');
insert into role (name) values ('ROLE_POTENTIAL_OWNER');
insert into role (name) values ('ROLE_CLIENT');

INSERT INTO admin (`id`, `email`, `first_name`, `last_name`, `password`, `username`, is_activated, pin, is_blocked, salt)
    VALUES (1, "health.care.clinic.psw+admin@gmail.com", "Admin", "Admirovic", "$2a$10$AY69sMB2v7eJHChMtGG61O4KPmkfj0sfPKKdp9vG1sMGsTqT4lbhq", "admin", true, 1111, false, "admin");

insert into admin_roles (admin_id, role_id) values (1, 1);

INSERT INTO company_owner (`id`,`blocked_date`,`email`,`first_name`,`forgotten`,`is_activated`,`is_blocked`,`last_name`,`missed_password_counter`,`password`,`pin`,`salt`,`username`)
    VALUES (1, null, "health.care.clinic.psw+owner@gmail.com", "Asd", 0, true, false, "ASD", 0, "$2a$10$S4KQX5hG/N3wvSXZ7ba/D.7Y7ELPLp4Q2AH.//3jJicoxnGd6AGwm",
    1234, "ownerr", "owner");


INSERT INTO company (`id`,`info`,`is_approved`,`name`,`company_owner_id`) VALUES (1, "dsf dfsA ASD ", true, "Kompanijica", 1);

INSERT INTO job_position (`id`,`avg_salary`,`name`,`company_id`) VALUES (1, 0, "Human resources", 1);
INSERT INTO job_position (`id`,`avg_salary`,`name`,`company_id`) VALUES (2, 0, "Psychologist", 1);

insert into owner_roles (company_owner_id, role_id) values (1, 2);

