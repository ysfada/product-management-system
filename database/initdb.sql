-- ADD YOUR INITIAL SQL HERE
-- INSERT INTO "public"."user" (
--     "username",
--     "password",
--     "is_active",
--     "is_staff",
--     "is_superuser"
-- )
-- VALUES
--     (
--         'admin',
--         '$argon2id$v=19$m=65536,t=1,p=2$7j31fNPcl9zEKSo3ADNVrQ$RWFWLoEGhUjUmM094e0/zvwiZ/HMWJeFvvKJaq7FX4s', -- 12345678
--         true,
--         true,
--         true
--     ),
--     (
--         'staff',
--         '$argon2id$v=19$m=65536,t=1,p=2$7j31fNPcl9zEKSo3ADNVrQ$RWFWLoEGhUjUmM094e0/zvwiZ/HMWJeFvvKJaq7FX4s', -- 12345678
--         true,
--         true,
--         false
--     ),
--     (
--         'user',
--         '$argon2id$v=19$m=65536,t=1,p=2$7j31fNPcl9zEKSo3ADNVrQ$RWFWLoEGhUjUmM094e0/zvwiZ/HMWJeFvvKJaq7FX4s', -- 12345678
--         true,
--         false,
--         false
--     )

-- INSERT INTO "public"."category" ("name", "description")
-- VALUES
--     ('phone', 'smart phones'),
--     ('camera', 'dslr cameras'),
--     ('television', 'televisions'),
--     ('clothing', 'clothings')

-- INSERT INTO "public"."product" ("name", "description", "category_id")
-- VALUES
--     ('IPhone X', '', 1),
--     ('LG G5', '', 1),
--     ('Samsung Note 6', '', 1),
--     ('Samsung Galaxy S21', '', 1),
--     ('Google Pixel 3', '', 1),
--     ('Canon EOS 90D', '', 2),
--     ('Canon EOS 7D Mark II', '', 2),
--     ('Samsung Smart Neo QLED TV', '', 3),
--     ('Sony Android Smart OLED TV', '', 3),
--     ('LG NanoCell Smart LED TV', '', 3),
--     ('Shirt', '', 4),
--     ('Coat', '', 4),
--     ('T-Shirt', '', 4),
--     ('Dress', '', 4)

-- INSERT INTO "public"."image"("name", "image_url", "thumbnail_url")
-- VALUES
--     (
--         'IPhone X',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'LG G5',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Samsung Note 6',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Samsung Galaxy S21',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Google Pixel 3',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Canon EOS 90D',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Canon EOS 7D Mark II',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Samsung Smart Neo QLED TV',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Sony Android Smart OLED TV',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'LG NanoCell Smart LED TV',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Shirt',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Coat',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'T-Shirt',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     ),
--     (
--         'Dress',
--         '/public/images/placeholder.png',
--         '/public/images/placeholder.png'
--     )

-- INSERT INTO "public"."product_images"("product_id", "image_id")
-- VALUES
--     (1, 1),
--     (2, 2),
--     (3, 3),
--     (4, 4),
--     (5, 5),
--     (6, 6),
--     (7, 7),
--     (8, 8),
--     (9, 9),
--     (10, 10),
--     (11, 11),
--     (12, 12),
--     (13, 13),
--     (14, 14)

-- INSERT INTO "public"."product_variant"(
--     "product_id",
--     "name",
--     "price",
--     "stock"
-- )
-- VALUES
--     (1, 'IPhone X', 11250.00, 2),
--     (1, 'IPhone X', 12350.00, 3),
--     (2, 'LG G5', 8450.00, 6),
--     (2, 'LG G5', 8830.00, 4),
--     (3, 'Samsung Note 6', 6500.00, 2),
--     (3, 'Samsung Note 6', 6450.00, 1),
--     (4, 'Samsung Galaxy S21', 9800.00, 10),
--     (4, 'Samsung Galaxy S21', 9750.00, 6),
--     (5, 'Google Pixel 3', 14500.00, 15),
--     (5, 'Google Pixel 3', 17500.00, 20),
--     (6, 'Canon EOS 90D', 22500.00, 3),
--     (6, 'Canon EOS 90D', 35750.00, 1),
--     (7, 'Canon EOS 7D Mark II', 45000.00, 2),
--     (7, 'Canon EOS 7D Mark II', 55000.00, 2),
--     (8, 'Samsung Smart Neo QLED TV', 18650.00, 3),
--     (8, 'Samsung Smart Neo QLED TV', 22800.00, 2),
--     (9, 'Sony Android Smart OLED TV', 18400.00, 2),
--     (9, 'Sony Android Smart OLED TV', 21550.00, 3),
--     (10, 'LG NanoCell Smart LED TV', 17650.00, 2),
--     (10, 'LG NanoCell Smart LED TV', 18550.00, 8),
--     (11, 'Shirt', 30.50, 30),
--     (11, 'Shirt', 36.00, 40),
--     (12, 'Coat', 600.00, 70),
--     (12, 'Coat', 450.00, 60),
--     (13, 'T-Shirt', 85.00, 100),
--     (13, 'T-Shirt', 80.00, 140),
--     (14, 'Dress', 675.00, 15),
--     (14, 'Dress', 500.00, 34)

-- INSERT INTO "public"."attribute"("name", "type")
-- VALUES
--     ('black', 'color'),
--     ('gold', 'color'),
--     ('silver', 'color'),
--     ('white', 'color'),
--     ('pink', 'color'),
--     ('ash', 'color'),
--     ('yellow', 'color'),
--     ('red', 'color'),
--     ('blue', 'color'),
--     ('green', 'color'),
--     (64, 'memory'),
--     (128, 'memory'),
--     (256, 'memory'),
--     (512, 'memory'),
--     ('tele', 'lens'),
--     ('macro', 'lens'),
--     ('black', 'lens'),
--     ('32"', 'size'),
--     ('40"', 'size'),
--     ('42"', 'size'),
--     ('56"', 'size'),
--     ('82"', 'size'),
--     ('XXL', 'size'),
--     ('XL', 'size'),
--     ('L', 'size'),
--     ('M', 'size'),
--     ('S', 'size'),
--     ('36', 'size'),
--     ('38', 'size'),
--     ('40', 'size'),
--     ('42', 'size'),
--     ('44', 'size')

-- INSERT INTO "public"."product_attributes"("product_variant_id", "attribute_id")
-- VALUES
--     (1, 1),
--     (1, 12),
--     (2, 2),
--     (2, 13),
--     (3, 1),
--     (3, 11),
--     (4, 1),
--     (4, 12),
--     (5, 3),
--     (5, 12),
--     (6, 3),
--     (6, 13),
--     (7, 3),
--     (7, 12),
--     (8, 4),
--     (8, 13),
--     (9, 3),
--     (9, 14),
--     (10, 2),
--     (10, 14),
--     (11, 15),
--     (12, 16),
--     (13, 15),
--     (14, 17),
--     (15, 18),
--     (16, 19),
--     (17, 20),
--     (18, 21),
--     (19, 19),
--     (20, 20),
--     (21, 1),
--     (21, 23),
--     (22, 2),
--     (22, 24),
--     (23, 3),
--     (23, 24),
--     (24, 4),
--     (24, 28),
--     (25, 2),
--     (25, 29),
--     (26, 1),
--     (26, 28),
--     (27, 7),
--     (27, 28),
--     (28, 8),
--     (28, 29)
