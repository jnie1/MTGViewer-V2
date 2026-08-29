-- deploy

ALTER TABLE transactions 
    ADD CONSTRAINT unique_exchange UNIQUE (log_group_id, from_container_id, to_container_id, scryfall_id),
    ADD CONSTRAINT positive_amount CHECK (amount > 0),
    ADD COLUMN log_group_id UUID;

UPDATE transactions
    SET log_group_id = group_id;

ALTER TABLE transactions
    ALTER COLUMN log_group_id SET NOT NULL,
    ADD CONSTRAINT log_group_fk
    FOREIGN KEY (log_group_id) REFERENCES log_groups(log_group_id)
    ON DELETE CASCADE,
    DROP COLUMN group_id,
    DROP COLUMN time;

-- rollback

ALTER TABLE transactions
    ADD COLUMN group_id UUID,
    ADD COLUMN time TIMESTAMP WITH TIME ZONE;

UPDATE transactions
    SET group_id = log_groups.log_group_id, time = log_groups.time
    FROM log_groups WHERE transactions.log_group_id = log_groups.log_group_id;

ALTER TABLE transactions 
    ALTER COLUMN group_id SET NOT NULL,
    ALTER COLUMN time SET NOT NULL,
    DROP CONSTRAINT unique_exchange,
    DROP CONSTRAINT positive_amount,
    DROP COLUMN log_group_id;