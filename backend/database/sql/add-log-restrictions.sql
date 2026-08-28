-- deploy

ALTER TABLE transactions 
    ADD CONSTRAINT unique_exchange UNIQUE (group_id, from_container_id, to_container_id, scryfall_id),
    ADD CONSTRAINT positive_amount CHECK (amount > 0);

-- rollback

ALTER TABLE transactions 
    DROP CONSTRAINT unique_exchange,
    DROP CONSTRAINT positive_amount;
