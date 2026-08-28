-- deploy

CREATE TABLE log_groups (
	log_group_id UUID PRIMARY KEY,
	time TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT
);

INSERT INTO log_groups (log_group_id, time)
SELECT DISTINCT group_id, time
FROM transactions;

-- rollback

DROP TABLE log_groups;