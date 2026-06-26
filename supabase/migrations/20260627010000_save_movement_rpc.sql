CREATE OR REPLACE FUNCTION save_movement(
    p_user_id UUID,
    p_real_distance_km FLOAT8,
    p_used_virtual_distance_km FLOAT8
) RETURNS TABLE (
    id UUID,
    user_id UUID,
    date DATE,
    real_distance_km FLOAT8,
    used_virtual_distance_km FLOAT8,
    created_at TIMESTAMPTZ
) AS $$
DECLARE
    v_today DATE := (NOW() AT TIME ZONE 'Asia/Tokyo')::DATE;
BEGIN
    UPDATE movements m
    SET
        real_distance_km = m.real_distance_km + p_real_distance_km,
        used_virtual_distance_km = m.used_virtual_distance_km + p_used_virtual_distance_km
    WHERE m.user_id = p_user_id
        AND m.date = v_today;

    IF NOT FOUND THEN
        INSERT INTO movements (user_id, date, real_distance_km, used_virtual_distance_km)
        VALUES (p_user_id, v_today, p_real_distance_km, p_used_virtual_distance_km);
    END IF;

    RETURN QUERY
    SELECT m.id, m.user_id, m.date, m.real_distance_km, m.used_virtual_distance_km, m.created_at
    FROM movements m
    WHERE m.user_id = p_user_id
        AND m.date = v_today;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
