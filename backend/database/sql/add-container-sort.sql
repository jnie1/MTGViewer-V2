-- deploy

ALTER TABLE containers
ADD sort_order INT NOT NULL DEFAULT 0

UPDATE containers
SET sort_order = container_id

-- rollback

ALTER TABLE containers
DROP column sort_order