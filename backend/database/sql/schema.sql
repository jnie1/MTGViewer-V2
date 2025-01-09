CREATE TABLE card (
	card_id SERIAL PRIMARY KEY,
	card_name VARCHAR(64) NOT NULL,
	scryfall_id UUID UNIQUE NOT NULL
);

CREATE TABLE container (
	container_id SERIAL PRIMARY KEY,
	container_name VARCHAR(32) NOT NULL,
	capacity INT NOT NULL,
	deletion_mark BOOLEAN NOT NULL
);

CREATE TABLE location (
	location_id SERIAL PRIMARY KEY,
	scryfall_id UUID NOT NULL,
	quantity INT NOT NULL,
	container_id INT NOT NULL,
	FOREIGN KEY (scryfall_id) REFERENCES card(scryfall_id)
	ON DELETE CASCADE,
	FOREIGN KEY (container_id) REFERENCES container(container_id)
	ON DELETE CASCADE
);

CREATE TABLE transaction (
	transaction_id SERIAL PRIMARY KEY,
	group_id INT NOT NULL,
	from_container INT,
	to_container INT,
	scryfall_id UUID NOT NULL,
	quantity INT NOT NULL,
	time TIMESTAMP WITH TIME ZONE NOT NULL,
	FOREIGN KEY (from_container) REFERENCES container(container_id),
	FOREIGN KEY (to_container) REFERENCES container(container_id)
);

CREATE TABLE users (
	user_id SERIAL PRIMARY KEY,
	email VARCHAR(128) NOT NULL,
	password_hash TEXT NOT NULL,
	role_string VARCHAR(32) NOT NULL
);