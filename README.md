# Product management

Product management system (WIP)

## Run Locally

Clone the project

```bash
git clone https://github.com/ysfada/product-management-system
```

Go to the project directory

```bash
cd product-management-system
```

Install dependencies

```bash
go get ./...
```

Read .env.example and create .env file.

Start postgres container

```bash
make up
```

Start the server

```bash
make run
```

## Deployment

To deploy this project run

```bash
  make build
```

## API Reference

<http://localhost:8080/docs/>

## Acknowledgements

- [Fiber - Express inspired web framework written in Go](https://github.com/gofiber/fiber)
- [Automatically generate RESTful API documentation with Swagger 2.0 for Go.](https://github.com/swaggo/swag)
- [The easiest way to create a README](https://readme.so/editor)
- [Create useful .gitignore files for your project](https://www.toptal.com/developers/gitignore)

## License

[MIT](https://choosealicense.com/licenses/mit/)
