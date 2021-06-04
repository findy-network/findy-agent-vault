UPDATE agent SET raw_jwt = '' WHERE raw_jwt IS NULL;
ALTER TABLE agent ALTER COLUMN raw_jwt SET NOT NULL;
UPDATE agent SET label = '' WHERE label IS NULL;
ALTER TABLE agent ALTER COLUMN label SET NOT NULL;


UPDATE connection SET approved = timestamp '0001-01-01' WHERE approved IS NULL;
ALTER TABLE connection ALTER COLUMN approved SET NOT NULL;
ALTER TABLE connection ALTER COLUMN approved SET DEFAULT timestamp '0001-01-01';

UPDATE connection SET archived = timestamp '0001-01-01' WHERE archived IS NULL;
ALTER TABLE connection ALTER COLUMN archived SET NOT NULL;
ALTER TABLE connection ALTER COLUMN archived SET DEFAULT timestamp '0001-01-01';
