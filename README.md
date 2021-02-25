# RCloud

## Build

Dependencies : go 1.16+ compiler, node/npm (Angular)

Install Angular CLI : `npm install -g @angular/cli`

Go to `internal/ui/gui`, and run `ng build --prod`.

Then go back to the root, and run `go build -tags prod ./cmd/rcloud`.

And your binary is here.

## Usage

```shell
./rcloud [-port=$port] [-https] [directory to manage] [path of the database]
```

The database will be created automatically at the given path, or if it exists it will load the content.

Note that each start of the application invalidates all sessions. You will have to log in after a restart.

`port` option let you choose on which port the server will be listening.
`https` option enables https mode. This does NOT start an https server. Instead, it enables security options like secure cookies. This option is useful to tell the application that it is running behind a ssl/tls reverse proxy.

## TODO
- DB is created relatively to the folder to list
- table style (at least on safari) when there is go back link
