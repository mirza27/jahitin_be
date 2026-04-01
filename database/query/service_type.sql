-- name: SeedServiceTypes :many
INSERT INTO
    service_types (name)
VALUES
    ('Jahit Baju'),
    ('Kecilkan'),
    ('Besarkan'),
    ('Pendekkan'),
    ('Resleting'),
    ('Ganti Kancing'),
    ('Tambal'),
    ('Tambah Furing'),
    ('Ganti Tali'),
    ('Ganti Ritsleting'),
    ('Pasang Bed') RETURNING *;

-- name: ListServiceTypes :many
SELECT
    *
FROM
    service_types
ORDER BY
    name ASC;