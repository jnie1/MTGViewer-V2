-- deploy

ALTER TABLE log_groups 
    ADD COLUMN description TEXT;

-- rollback

ALTER TABLE log_groups
    DROP COLUMN description;
