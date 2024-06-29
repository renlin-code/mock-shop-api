# Mock Shop API


Mock Shop API is an API service for a fake online shop. It allows you to create a catalog for products and product categories. It also has authorization for customers where they can place orders and view their purchase history.
![Texto alternativo](https://imgur.com/JaTTmzG.png)

## Building and running the App

This app uses Docker, which makes building easier. You just need to install [Docker](https://docs.docker.com/) on your system.

In addition to this you need to set your environment variables in a .env file. As an example see the [.env.example](https://github.com/renlin-code/mock-shop-api/blob/master/.env.example) file


Having Docker on your system and environment variables setted, you only need to run:

```bash
docker-compose up
```

## Documentation

After building and running the app you can find the Documentation at the address "/swagger/index.html"
