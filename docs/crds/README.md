# CRD documentation

The CRDs documentation is generated automatically by using <https://github.com/elastic/crd-ref-docs> and the `config.yml` file shipped within this directory.

## Documentation generation

### Markdown

To generate markdown documentation:

```shell
make markdown
```

The result will be save to the `crd.md` file.

### ASCIIDoc

To generate asciidoc documentation:

```shell
make asciidoc
```

The result will be save to the `crd.adoc` file.

## Development notes

Ensure the contents of the `templates` directory match the version of
`crd-ref-docs` being used.
