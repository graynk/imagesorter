# imagesorter

Yet another very specific tool that only I need. By default it requires a terminal that supports [terminal graphics protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/),
but you can fallback to Sixel by passing `--sixel` to the command.

https://github.com/graynk/imagesorter/assets/3626328/7535af34-4030-4955-ad69-e5c4b9806156

You can install it by grabbing a binary from [releases](https://github.com/graynk/imagesorter/releases) page, or by running
```bash
go install github.com/graynk/imagesorter@latest
```

You use it like this:
```bash
imagesorter [--sixel] /path/to/source cool_pictures not_so_cool_pictures can_be_deleted
```
or like this
```bash
imagesorter [--sixel] --scan /path/to/source /path/to/target
```

Then it will read every PNG/JPG in source directory, display it in the terminal and ask you to which of the target directories it should be moved.
Then it will move it accordingly.

Optional `--scan` argument will scan every directory it `target` and create a list of target directories that way.

Note: it uses `os.Rename` to move the file because I am lazy, so it won't move the file between different drives.
