**krab** is a tiny utility for extracting `mergedimage.png` from .kra files (Krita project files), which avoids the slowness of [Krita's command-line exporting](https://docs.krita.org/en/reference_manual/linux_command_line.html).

# Install

```bash
go install github.com/greenthepear/krab@main
```

# Usage

Check with
```bash
krab -h
```

```bash
krab -o ../Krita/card_pattern_blue.png card_pattern_blue.kra
```

You can also omit `-o` and it will create a file in the same directory with `.kra` changed to `.png`.

```bash
krab card_pattern_blue.kra
```

If you specify a directory it will output the png there
```bash
krab -o Graphics/ card_pattern_blue.kra
```

Globs are supported
```bash
krab *.kra
```
