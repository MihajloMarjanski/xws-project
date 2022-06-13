insert into role (name) values ('ROLE_ADMIN');
insert into role (name) values ('ROLE_COMPANY_OWNER');
insert into role (name) values ('ROLE_POTENTIAL_OWNER');
insert into role (name) values ('ROLE_CLIENT');

insert into permission (id, name) values (1, 'approveCompany');
insert into permission (id, name) values (2, 'adminByUsername');
insert into permission (id, name) values (3, 'updateAdmin');
insert into permission (id, name) values (4, 'updateClient');
insert into permission (id, name) values (5, 'clientByUsername');
insert into permission (id, name) values (6, 'updateCompanyOwner');
insert into permission (id, name) values (7, 'createCompany');
insert into permission (id, name) values (8, 'getOwner');
insert into permission (id, name) values (9, 'createComment');
insert into permission (id, name) values (10, 'updateSalary');
insert into permission (id, name) values (11, 'addInformation');
insert into permission (id, name) values (12, 'createJobOffer');
insert into permission (id, name) values (13, 'ownerByUsername');
insert into permission (id, name) values (14, 'apiKey');

insert into role_permissions (role_id, permission_id) values (1, 1);
insert into role_permissions (role_id, permission_id) values (1, 2);
insert into role_permissions (role_id, permission_id) values (1, 3);
insert into role_permissions (role_id, permission_id) values (4, 4);
insert into role_permissions (role_id, permission_id) values (4, 5);
insert into role_permissions (role_id, permission_id) values (2, 6);
insert into role_permissions (role_id, permission_id) values (3, 6);
insert into role_permissions (role_id, permission_id) values (3, 7);
insert into role_permissions (role_id, permission_id) values (2, 8);
insert into role_permissions (role_id, permission_id) values (4, 9);
insert into role_permissions (role_id, permission_id) values (4, 10);
insert into role_permissions (role_id, permission_id) values (4, 11);
insert into role_permissions (role_id, permission_id) values (2, 13);
insert into role_permissions (role_id, permission_id) values (3, 13);
insert into role_permissions (role_id, permission_id) values (2, 14);


INSERT INTO admin (`id`, `email`, `first_name`, `last_name`, `password`, `username`, is_activated, pin, is_blocked, salt)
    VALUES (1, "health.care.clinic.psw+admin@gmail.com", "Admin", "Admirovic", "$2a$10$AY69sMB2v7eJHChMtGG61O4KPmkfj0sfPKKdp9vG1sMGsTqT4lbhq", "admin", true, "1111", false, "admin");

insert into admin_roles (admin_id, role_id) values (1, 1);

INSERT INTO company_owner (`id`,`blocked_date`,`email`,`first_name`,`forgotten`,`is_activated`,`is_blocked`,`last_name`,`missed_password_counter`,`password`,`pin`,`salt`,`username`)
    VALUES (1, null, "owner@gmail.com", "Asd", 0, true, false, "ASD", 0, "$2a$10$DRq8KussSpZpBw57z3mpsOBuBUI7D1Ut/tR5PWef5XebPdm1.VCL.",
            "$2a$10$b9vQ1SxHQX7z06qJd7Jgt.22WwJ9PP81ERmii1nBwDCTYs2NEI0Mm", "QWEQWE", "owner");


INSERT INTO company (`id`,`info`,`is_approved`,`name`,`company_owner_id`) VALUES (1, "dsf dfsA ASD ", true, "Kompanijica", 1);

INSERT INTO job_position (`id`,`avg_salary`,`name`,`company_id`) VALUES (1, 0, "Human resources", 1);
INSERT INTO job_position (`id`,`avg_salary`,`name`,`company_id`) VALUES (2, 0, "Psychologist", 1);

insert into owner_roles (company_owner_id, role_id) values (1, 2);


INSERT INTO client (`id`,`blocked_date`,`email`,`first_name`,`forgotten`,`is_activated`,`is_blocked`,`last_name`,`missed_password_counter`,`password`,`pin`,`salt`,`username`)
    VALUES (1, null, "asd@gmail.com", "Mile", 0, true, false, "Kitic", 0, "$2a$10$fIYdk9bNYb2QzrKJk5kNIeSSyFFKCnUnsVGsnJBRIpNZ0XLMtJ4h2",
            "$2a$10$l9U3nJjCqQTyBfIUxko63.NfQvd6zLrd6o6zfZ5Aa.K32nNutdlrW", "ASDASD", "mile" );

insert into client_roles (client_id, role_id) values (1, 4);


