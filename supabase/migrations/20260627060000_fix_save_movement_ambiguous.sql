DROP FUNCTION IF EXISTS save_movement(UUID, FLOAT8, FLOAT8);

CREATE FUNCTION save_movement(
    p_user_id UUID,
    p_real_distance_km FLOAT8,
    p_used_virtual_distance_km FLOAT8
) RETURNS SETOF movements
LANGUAGE plpgsql SECURITY DEFINER
SET search_path = public
AS $$
DECLARE
    v_today DATE := (NOW() AT TIME ZONE 'Asia/Tokyo')::DATE;
BEGIN
    INSERT INTO movements (user_id, date, real_distance_km, used_virtual_distance_km)
    VALUES (p_user_id, v_today, p_real_distance_km, p_used_virtual_distance_km)
    ON CONFLICT ON CONSTRAINT movements_user_id_date_key DO UPDATE
    SET
        real_distance_km         = movements.real_distance_km + p_real_distance_km,
        used_virtual_distance_km = movements.used_virtual_distance_km + p_used_virtual_distance_km;

    RETURN QUERY
    SELECT * FROM movements
    WHERE movements.user_id = p_user_id
      AND movements.date    = v_today;
END;
$$;
