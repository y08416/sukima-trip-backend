CREATE OR REPLACE FUNCTION get_total_distance(p_user_id uuid)
RETURNS float8
LANGUAGE sql
SECURITY DEFINER
SET search_path = public
AS $$
  SELECT COALESCE(SUM(real_distance_km), 0)
  FROM movements
  WHERE user_id = p_user_id;
$$;
