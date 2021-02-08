# RCloud

## Build

Dependencies : go 1.16+ compiler, node/npm (Angular)

Install Angular CLI : `npm install -g @angular/cli`

Go to `internal/ui/gui`, and run `ng build`.

Then go back to the root, and run `go build -tags prod ./cmd/rcloud`.

And your binary is here.

## Usage

```shell
./rcloud [-port=$port] [directory to manage] [path of the database]
```

The database will be created automatically at the given path, or if it exists it will load the content.

Note that each start of the application invalidates all sessions. You will have to log in after a restart.