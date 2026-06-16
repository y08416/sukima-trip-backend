alter table visited_places
  add constraint visited_places_user_id_place_id_key unique (user_id, place_id);
