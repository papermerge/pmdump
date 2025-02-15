# PMG Dump

Tool to migrate Papermerge data from one version to another.
Use this tool to migrate from Papermerge 2.0, 2.1, 3.3 to the latest version e.g 3.4.

Work in progress...

## Config

Work in progress...

Example of config.yaml file

```yaml
media_root: /path/to/media/folder/
database_url: /path/to/data/papermerge.db
target_file: "mydb.tar.gz"
```

## Usage

Work in progress...

```
$ pmg_dump -c export.yaml export
```

```
$ pmg_dump -c import.yaml import -f /path/to/archive.tar.gz
```
