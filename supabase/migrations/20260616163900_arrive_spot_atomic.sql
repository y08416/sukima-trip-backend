create or replace function arrive_spot(
  p_user_id    uuid,
  p_place_id   text,
  p_place_name text,
  p_amount     int
)
returns int
language plpgsql
security definer
set search_path = public
as $$
declare
  v_balance int;
begin
  insert into visited_places (user_id, place_id, place_name)
  values (p_user_id, p_place_id, p_place_name);

  update coins
  set balance = balance + p_amount
  where user_id = p_user_id;

  if not found then
    raise exception 'ユーザーのコインレコードが見つかりません: %', p_user_id;
  end if;

  select balance into v_balance
  from coins
  where user_id = p_user_id;

  return v_balance;
end;
$$;
