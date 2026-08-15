-- deploy

ALTER TABLE card_deposits
    ADD CONSTRAINT unique_container_card UNIQUE (scryfall_id, container_id),
    ADD CONSTRAINT positive_amount CHECK (amount > 0),
    ADD COLUMN oracle_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;

-- rollback

ALTER TABLE card_deposits
    DROP CONSTRAINT unique_container_card,
    DROP CONSTRAINT positive_amount,
    DROP COLUMN oracle_id;