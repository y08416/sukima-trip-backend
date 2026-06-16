ALTER TABLE spot_likes
  ADD CONSTRAINT spot_likes_user_id_place_id_key UNIQUE (user_id, place_id);
