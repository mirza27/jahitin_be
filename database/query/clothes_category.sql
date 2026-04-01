-- name: SeedClothesCategories :many
INSERT INTO
    clothes_categories (name)
VALUES
    ('Kaos'),
    ('Kebaya'),
    ('Baju Anak'),
    ('Jas'),
    ('Jaket'),
    ('Seragam'),
    ('Kemeja'),
    ('Celana'),
    ('Rok'),
    ('Gaun'),
    ('Gamis') RETURNING *;

-- name: ListClothesCategories :many
SELECT
    *
FROM
    clothes_categories
ORDER BY
    name ASC;