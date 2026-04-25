CREATE OR REPLACE FUNCTION infer_age_group(birth_date DATE) RETURNS TEXT AS $$
    SELECT CASE
        WHEN EXTRACT(YEAR FROM age(birth_date)) < 1  THEN 'infant'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 1  AND 2  THEN 'toddler'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 3  AND 5  THEN 'pre_school'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 6  AND 11 THEN 'middle_childhood'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 12 AND 14 THEN 'young_teen'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 15 AND 17 THEN 'teenager'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 18 AND 24 THEN 'young_adult'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 25 AND 39 THEN 'early_adult'
        WHEN EXTRACT(YEAR FROM age(birth_date)) BETWEEN 40 AND 59 THEN 'middle_aged_adult'
        WHEN EXTRACT(YEAR FROM age(birth_date)) >= 60              THEN 'old_adult'
        ELSE 'unknown'
    END;
$$ LANGUAGE sql IMMUTABLE;
