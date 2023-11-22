# Dagger CI Tool

This tool wraps any command into a Go container to guarantee a clean environment.
Any artifacts stored at 'output/' will be exported from the container.
Example usage:

```console
ci -cmd "go build -o output/"
```
