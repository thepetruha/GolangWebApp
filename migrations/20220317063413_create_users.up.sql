CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null
);

CREATE TABLE todo (
    id_td bigserial not null primary key,
    customer_id int not null,
    text_todo varchar not null,
    date_todo date not null,

   	FOREIGN KEY(customer_id) REFERENCES users(id) ON DELETE CASCADE
);