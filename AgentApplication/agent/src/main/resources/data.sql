insert into role (name) values ('ROLE_ADMIN');
insert into role (name) values ('ROLE_COMPANY_OWNER');
insert into role (name) values ('ROLE_POTENTIAL_OWNER');
insert into role (name) values ('ROLE_REGISTERED_USER');

INSERT INTO admin (`id`, `email`, `first_name`, `last_name`, `password`, `username`, is_activated, pin)
    VALUES (1, "admin@gmail.com", "Admin", "Admirovic", "$2a$10$1P.3BtNc4h5aC7ZTDUhM6OM9/kYw5jkalw0cIDEtWLQqaTGPuMXju", "admin", true, 1111);

insert into admin_roles (admin_id, role_id) values (1, 1);
