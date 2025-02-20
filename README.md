# PMDump

Tool to migrate Papermerge data from one version to another.
Use this tool to migrate from Papermerge 2.0, 2.1, 3.3 to the latest version e.g 3.4.

Work in progress...

## Usage


### Export

To export data **from Papermerge 2.0**:

create `source.yaml` file

```yaml
version: 2.0
media_root: /path/to/media/folder/
database_url: /path/to/data/papermerge.db
```

Run following command:

```
$ pmdump -c source.yaml -f pm2.0.tar.gz export
```

It will create `pm2.0.tar.gz` archive.


### Import


```
$ pmg_dump -c import.yaml -f /path/to/archive.tar.gz import
```

Note that `export` or `import` commands are at the end of parameters list.
