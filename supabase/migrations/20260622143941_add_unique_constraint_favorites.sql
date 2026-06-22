ALTER TABLE favorites
  ADD CONSTRAINT favorites_user_id_place_id_key UNIQUE (user_id, place_id);
