
CREATE TABLE products (
    product_id INT PRIMARY KEY,
    product_name TEXT NOT NULL,
    price FLOAT NOT NULL,
    stock INT NOT NULL
);

CREATE TABLE items (
    item_id INT PRIMARY KEY,
    item_name TEXT NOT NULL,
    price FLOAT NOT NULL,
    quantity INT NOT NULL
);
