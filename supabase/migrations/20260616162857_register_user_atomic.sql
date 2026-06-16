create or replace function register_user(
  p_user_id uuid,
  p_name    text,
  p_gender  text
)
returns void
language plpgsql
security definer
set search_path = public
as $$
begin
  insert into users (id, name, gender)
  values (p_user_id, p_name, p_gender);

  insert into coins (user_id, balance)
  values (p_user_id, 0);
end;
$$;
