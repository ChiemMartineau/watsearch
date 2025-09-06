# WATsearch

A tool to navigate University of Waterloo websites!

## Setup Development Environment

1. Start docker containers

   ```bash
   docker compose -f ./docker/docker-compose.yml up -d
   ```

2. Apply database migrations

   ```bash
   goose up
   ```

3. You can now use the WATsearch tools. You should start with running the `scraper`.

## Database Migrations

Apply migrations:

```bash
goose up
```

Roolback migrations by 1 version:

```bash
goose down
```

For more information, see the [Goose Documentation](https://pressly.github.io/goose/documentation/cli-commands/)
