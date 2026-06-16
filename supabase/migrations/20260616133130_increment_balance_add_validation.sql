CREATE OR REPLACE FUNCTION increment_balance(p_user_id UUID, p_amount INT)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_amount <= 0 THEN
    RAISE EXCEPTION 'p_amount must be positive, got %', p_amount;
  END IF;

  UPDATE coins
  SET balance = balance + p_amount
  WHERE user_id = p_user_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'user % not found in coins', p_user_id;
  END IF;
END;
$$;
