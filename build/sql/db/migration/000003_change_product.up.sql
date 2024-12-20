ALTER TABLE product ADD product_count int;
UPDATE product SET product_count = 0;
ALTER TABLE product ADD total int;
UPDATE product SET total = 0;