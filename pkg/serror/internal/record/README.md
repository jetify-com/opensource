# Package record

This is a direct copy of the corresponding code from the `slog` package with
a few slight modifications we need for serror:

- Add support for Int values
- Remove the need for Level
- Add utility functions to convert to slog.Record (used when logging)