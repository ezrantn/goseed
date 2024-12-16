# goseed - database seeder

goseed is a library to seed random data into a PostgreSQL database using models and `faker`-generated values for testing, development, or prototyping.

## Installation

To install goseed, run:

```bash
go get github.com/ezrantn/goseed
```

## Features

- **Easy to Use:** Seed your database with random data by defining models and specifying the number of rows to generate.
- **Faker Integration:** Integrates with the faker library to generate realistic random data like names, emails, numbers, etc.
- **Customizable:** Define your own table models with custom fields and data types.
- **Supports Multiple Tables:** Seed multiple tables in a single run.
  
## Usage

See the `examples/` directory for implementation examples.

For the `faker` struct tag, refer to their documentation [here](https://github.com/go-faker/faker/blob/main/example_with_tags_test.go).

## License

This tool is open-source and available under the [BSD-3 Clause](https://github.com/ezrantn/goseed/blob/main/LICENSE) License.

## Contributions

Contributions are welcome! Please feel free to submit a pull request.
