do $$
begin
  if not exists (
    select 1 from pg_constraint
    where conname = 'visited_places_user_id_place_id_key'
  ) then
    alter table visited_places
      add constraint visited_places_user_id_place_id_key unique (user_id, place_id);
  end if;
end $$;
