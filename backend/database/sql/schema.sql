CREATE TABLE containers (
	container_id SERIAL PRIMARY KEY,
	container_name VARCHAR(32) NOT NULL,
	capacity INT NOT NULL,
	deletion_mark BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE card_deposits (
	card_deposit_id SERIAL PRIMARY KEY,
	scryfall_id UUID UNIQUE NOT NULL,
	amount INT NOT NULL,
	container_id INT NOT NULL,
	FOREIGN KEY (container_id) REFERENCES containers(container_id)
	ON DELETE CASCADE
);

CREATE TABLE transactions (
	transaction_id SERIAL PRIMARY KEY,
	group_id UUID NOT NULL,
	from_container_id INT,
	to_container_id INT,
	scryfall_id UUID NOT NULL,
	amount INT NOT NULL,
	time TIMESTAMP WITH TIME ZONE NOT NULL,
	FOREIGN KEY (from_container_id) REFERENCES containers(container_id),
	FOREIGN KEY (to_container_id) REFERENCES containers(container_id)
);

CREATE TABLE users (
	user_id SERIAL PRIMARY KEY,
	email VARCHAR(128) NOT NULL,
	password_hash TEXT NOT NULL,
	role_string VARCHAR(32) NOT NULL
);