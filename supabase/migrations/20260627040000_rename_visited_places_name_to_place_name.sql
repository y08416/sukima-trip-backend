ALTER TABLE visited_places DROP COLUMN IF EXISTS place_name;
ALTER TABLE visited_places RENAME COLUMN name TO place_name;
