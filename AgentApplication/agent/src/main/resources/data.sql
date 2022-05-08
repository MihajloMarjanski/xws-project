insert into role (name) values ('ROLE_ADMIN');
insert into role (name) values ('ROLE_COMPANY_OWNER');
insert into role (name) values ('ROLE_REGISTERED');

INSERT INTO admin (`id`, `email`, `first_name`, `last_name`, `password`, `username`, `role_id`)
    VALUES (1, "admin@gmail.com", "Admin", "Admirovic", "$2a$10$1P.3BtNc4h5aC7ZTDUhM6OM9/kYw5jkalw0cIDEtWLQqaTGPuMXju", "admin", 1);
