**krab** is a tiny utility for extracting `mergedimage.png` from .kra files (Krita project files), which avoids the slowness of [Krita's command-line exporting](https://docs.krita.org/en/reference_manual/linux_command_line.html).

# Install

```bash
go install module github.com/greenthepear/krab@main
```

# Usage

```bash
krab -i card_pattern_blue.kra -o ../Krita/card_pattern_blue.png
```

You can also omit -o and it will create a file in the same directory with `.kra` changed to `.png`.