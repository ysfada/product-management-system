INSERT INTO "public"."user" (
    "username",
    "password",
    "is_active",
    "is_staff",
    "is_superuser"
)
VALUES
    (
        'admin',
        '$argon2id$v=19$m=65536,t=1,p=2$7j31fNPcl9zEKSo3ADNVrQ$RWFWLoEGhUjUmM094e0/zvwiZ/HMWJeFvvKJaq7FX4s', -- 12345678
        true,
        true,
        true
    ),
    (
        'staff',
        '$argon2id$v=19$m=65536,t=1,p=2$7j31fNPcl9zEKSo3ADNVrQ$RWFWLoEGhUjUmM094e0/zvwiZ/HMWJeFvvKJaq7FX4s', -- 12345678
        true,
        true,
        false
    ),
    (
        'user',
        '$argon2id$v=19$m=65536,t=1,p=2$7j31fNPcl9zEKSo3ADNVrQ$RWFWLoEGhUjUmM094e0/zvwiZ/HMWJeFvvKJaq7FX4s', -- 12345678
        true,
        false,
        false
    )

