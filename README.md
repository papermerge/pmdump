# PMDump

Tool to migrate Papermerge data from one version to another.
Use this tool to migrate from Papermerge 2.0, 2.1, 3.3 to the latest version e.g 3.4.

Work in progress...

## Usage


### Export Config

create `source.yaml` file in which you need specify
version of the Papermerge DMS from where to export, media
root folder, and database.

For example, to export from Papermerge DMS 2.0 (i.e. app version 2.0),
which stores its data in SQLite database:

```yaml
app_version: 2.0
media_root: /path/to/media/folder/
database_url: /path/to/data/papermerge.db
```

In case of exporting data from app version 3.3 which
store data in postgres database, configuration yaml file
would look like:

```yaml
media_root: /home/eugen/DockerCompose/pm3.3.1-pg/media
database_url: postgresql://coco:jumbo@127.0.0.1:5432/pmgdb?sslmode=disable
app_version: 3.3
```

`app_version` - it the version of Papermerge DMS from which you intend to export data. Version must be specified in `<major>.<minor>` format. Here are some examples of valid values: `2.0`, `2.1`, `3.0`, `3.2`, `3.4`.

`media_root` is folder where Papermerge DMS stores actual files. `pmdump` should have access to this folder.
Papermerge DMS uses `PAPERMERGE__MAIN__MEDIA_ROOT` environment variable to configure its media root folder. If not configured, it defaults to `/core_app/media/`.

`database_url` - is the database URL. For SQLite it is the path to sqlite3 file, prefixed with `sqlite://`. Some examples: `sqlite:///home/eugen/DockerCompose/pm3.4-sqlite/data/db.sqlite3`, `sqlite:///var/data/db.sqlite3`. Notice in examples
there are three `///`: two stand for the sqlite scheme, and one stands for root folder.

For PostgreSQL scheme must either `postgres://` or `postgresql://`. Here are some valid examples: `postgresql://coco:jumbo@127.0.0.1:5432/pmgdb?sslmode=disable`, `postgresql://coco:jumbo@127.0.0.1:5432/pmgdb`

### Export Command

Export command will save all data into zipped tar archive. You specify name of the exported file via `-f` flag and configuration file via `-c` flag:

```
$ pmdump -c source.yaml -f pm3.3.tar.gz export
```

Note `export` command **must be at the end of parameters list**.
Placing `export` command BEFORE parameters list won't work:

```
// wrong! "export" command must be at the end!
$ pmdump export -c source.yaml -f pm3.3.tar.gz
```

If you place `export` before parameters flag, you will get "missing configuration" error.

The general form for export command is:

```
$ pmdump <flag parameters> export
```

Here are couple of valid examples for `export` command:

```
$ pmdump -c examples/source_3_3_pg.yaml -f pm3.3_pg.tar.gz export
```



### Import


Create `dest.yaml` file

```yaml
version: 3.4
media_root: /path/to/media/folder/
database_url: /path/to/data/papermerge.db
```

```
$ pmdump -c dest.yaml -f /path/to/archive.tar.gz import
```

Note that `export` or `import` commands are at the end of parameters list.


## Development

You need golang version 1.24:

```
$ go version

go version go1.24.0
```

To compile from source code use following commands:

```
$ go mod download
$ go mod tidy
$ go build dump.go
```

Above commands will create `pmdump` binary.
