CREATE TABLE users (
    id SERIAL NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    profile_image VARCHAR(255)
);

CREATE TABLE categories (
    id SERIAL NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255),
    image_url VARCHAR(255)
);

CREATE TABLE orders (
    id SERIAL NOT NULL UNIQUE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    date DATE NOT NULL
);

CREATE TABLE products (
    id SERIAL NOT NULL UNIQUE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    price FLOAT NOT NULL,
    sale_price FLOAT,
    images_url TEXT,
    available BOOLEAN NOT NULL,
    stock INT NOT NULL
);

CREATE TABLE orders_products (
    id SERIAL NOT NULL UNIQUE,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE NOT NULL,
    product_id INT REFERENCES products(id) ON DELETE CASCADE NOT NULL
);
