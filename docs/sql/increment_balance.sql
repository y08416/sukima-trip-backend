-- Supabase SQL Editor で実行する
-- コイン残高をアトミックにインクリメントする関数
CREATE OR REPLACE FUNCTION increment_balance(p_user_id UUID, p_amount INT)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE coins
  SET balance = balance + p_amount
  WHERE user_id = p_user_id;
END;
$$;
