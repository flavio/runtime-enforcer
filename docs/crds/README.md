# CRD documentation

The CRDs documentation is generated automatically by using <https://github.com/elastic/crd-ref-docs> and the `config.yml` file shipped within this directory.

## Documentation generation

To generate the AsciiDoc documentation:

```shell
make asciidoc
```

The result will be saved to the `CRD-docs-for-docs-repo.adoc` file.

## Development notes

Ensure the contents of the `templates` directory match the version of
`crd-ref-docs` being used.
