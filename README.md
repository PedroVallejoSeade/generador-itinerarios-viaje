# generador-itinerarios-viaje

## City Search Spike

`citysearch` is a terminal tool that looks up cities by name and lists matching
cities with disambiguating country/region context, ranked by population and
capped at 10 results. The city dataset is bundled into the binary, so the tool
works fully offline with no API key or account.

### Build

```bash
go build -o bin/citysearch ./cmd/citysearch
# or
make build
```

### Usage

```bash
citysearch <city-name>
citysearch -h | --help
```

The city name is a positional argument; quote multi-word names
(e.g. `citysearch "san jose"`).

### Examples

```bash
$ ./bin/citysearch paris
Paris, Île-de-France, France
...

$ ./bin/citysearch springfield
Springfield, Missouri, United States
Springfield, Massachusetts, United States
Springfield, Illinois, United States
...

$ ./bin/citysearch zzzzzz
No cities found matching "zzzzzz".
```

Output is one match per line as `"<Name>, <Region>, <Country>"` (region omitted
when unknown), ordered by population (largest first) and capped at 10 results.

### Exit codes

| Code | Meaning |
|------|---------|
| 0 | Query processed successfully (including a valid "no results" case) |
| 1 | Invalid usage — empty/whitespace query or missing argument |
| 2 | Data-source/load error |

### Tests

```bash
go test ./...
# or
make check   # gofmt + go vet + go test
```

## Data attribution

City data © [GeoNames](https://www.geonames.org/) (`cities15000` dump),
licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). The
dump is converted into a `city,country,admin_name,population` CSV and embedded in
the binary via `//go:embed`. See
[specs/001-city-search-spike/research.md](specs/001-city-search-spike/research.md)
for the data-source evaluation and rationale.