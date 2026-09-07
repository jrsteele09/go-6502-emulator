# Klaus Dormann 6502 Functional Tests

Source: https://github.com/Klaus2m5/6502_65C02_functional_tests

The `TestKlausFunctional` and `TestKlausDecimal` integration tests download the
upstream GPL-3.0 test artifacts into this directory when they are not already
cached. The downloaded files are intentionally ignored so this MIT-licensed
repository does not vendor GPL test assets by default.

Run the tests with:

```bash
go test ./cpu -run Klaus
```
