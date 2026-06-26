alter table visited_places add column coin_amount integer not null default 10;

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
  if p_amount <= 0 then
    raise exception 'コイン付与枚数は1以上である必要があります';
  end if;

  insert into visited_places (user_id, place_id, place_name, coin_amount)
  values (p_user_id, p_place_id, p_place_name, p_amount);

  update coins
  set balance = balance + p_amount
  where user_id = p_user_id
  returning balance into v_balance;

  if not found then
    raise exception 'ユーザーのコインレコードが見つかりません: %', p_user_id;
  end if;

  return v_balance;
end;
$$;
