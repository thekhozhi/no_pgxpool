CREATE TYPE user_type_enum AS ENUM ('admin', 'customer');

CREATE TABLE users (
  id UUID PRIMARY KEY,
  full_name VARCHAR(30),
  phone VARCHAR(13),
  password VARCHAR(128),
  cash INT DEFAULT 0,
  user_type user_type_enum
);

CREATE TABLE categories (
  id UUID PRIMARY KEY,
  name VARCHAR(30)
);

CREATE TABLE products (
  id UUID PRIMARY KEY,
  name VARCHAR(30),
  price INT,
  original_price INT,
  quantity INT,
  category_id UUID REFERENCES categories(id)
);

CREATE TABLE baskets (
  id UUID PRIMARY KEY,
  customer_id UUID REFERENCES users(id),
  total_sum INT
);

CREATE TABLE basket_products (
  id UUID PRIMARY KEY,
  basket_id UUID REFERENCES baskets(id),
  product_id UUID REFERENCES products(id),
  quantity INT
);


INSERT INTO users (id, full_name, phone, password, cash, user_type)
VALUES 
  ('c6ecf5e6-8d3a-4b5b-8df8-2c58e0e6a9f1', 'John Doe', '1234567890', 'password123', 100, 'customer'),
  ('84d32d47-15e1-4a0c-ba0a-9c3b5e7d5aaf', 'Jane Smith', '9876543210', 'securepass', 50, 'customer');


INSERT INTO categories (id, name)
VALUES 
  ('6c9e27b8-d66f-4d2c-8c1d-3c6c9a6d7e0c', 'Electronics'),
  ('b8b7a40f-2f00-4c0d-8d46-4a3d5e6c7b8a', 'Clothing');


INSERT INTO products (id, name, price, original_price, quantity, category_id)
VALUES 
  ('b07f4e72-4f8e-4b79-9f8c-4a5e6d7f8b09', 'Smartphone', 500, 700, 10, '6c9e27b8-d66f-4d2c-8c1d-3c6c9a6d7e0c'),
  ('a94b7cd3-0ead-4b7e-8d0f-3c2b1a9e8d76', 'T-Shirt', 20, 30, 20, 'b8b7a40f-2f00-4c0d-8d46-4a3d5e6c7b8a');

 
INSERT INTO baskets (id, customer_id, total_sum)
VALUES 
  ('7d891c50-1e97-4d05-8e4f-0d1c2b3a4b5c', 'c6ecf5e6-8d3a-4b5b-8df8-2c58e0e6a9f1', 700),
  ('d8e9f0a1-b2c3-4d5e-6f7a-8b9c0d1e2f3a', '84d32d47-15e1-4a0c-ba0a-9c3b5e7d5aaf', 50);

 
INSERT INTO basket_products (id, basket_id, product_id, quantity)
VALUES 
  ('f1e2d3c4-b5a6-4c7d-8e9f-0a1b2c3d4e5f', '7d891c50-1e97-4d05-8e4f-0d1c2b3a4b5c', 'b07f4e72-4f8e-4b79-9f8c-4a5e6d7f8b09', 2),
  ('9a8b7c6d-5e4f-4d3c-2b1a-0f9e8d7c6b5a', 'd8e9f0a1-b2c3-4d5e-6f7a-8b9c0d1e2f3a', 'a94b7cd3-0ead-4b7e-8d0f-3c2b1a9e8d76', 1);
