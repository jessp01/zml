# zml

[![CI][badge-build]][build]
[![GoDoc][go-docs-badge]][go-docs]
[![GoReportCard][go-report-card-badge]][go-report-card]
[![License][badge-license]][license]

Diagram and flowchart tool written in Golang

### Installation

```sh
$ go install github.com/jessp01/zml/cmd/zml_cli@latest
```

### Example

See the [examples dir](./examples) for sample input files.

```sh
$ zml_cli ./examples/sequence_flow1.zml
```

Will create `./examples/sequence_flow1.zml.png`

[![seq chart](examples/sequence_flow1.zml.png)]

[license]: ./LICENSE
[badge-license]: https://img.shields.io/github/license/jessp01/zml.svg
[go-docs-badge]: https://godoc.org/github.com/jessp01/zml?status.svg
[go-docs]: https://godoc.org/github.com/jessp01/zml
[go-report-card-badge]: https://goreportcard.com/badge/github.com/jessp01/zml
[go-report-card]: https://goreportcard.com/report/github.com/jessp01/zml
[badge-build]: https://github.com/jessp01/zml/actions/workflows/go.yml/badge.svg
[build]: https://github.com/jessp01/zml/actions/workflows/go.yml

