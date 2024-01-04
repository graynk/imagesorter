# imagesorter

Yet another very specific tool that only I need. It requires a terminal that supports [terminal graphics protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/)

https://github.com/graynk/imagesorter/assets/3626328/e0d80d12-c202-417a-a318-6f64fde4de6d

You use it like this:
```bash
imagesorter /path/to/source cool_pictures not_so_cool_pictures can_be_deleted
```

Then it will read every PNG/JPG in source directory, display it in the terminal and ask you to which of the target directories it should be moved.
Then it will move it accordingly. 

Note: it uses `os.Rename` to move the file because I am lazy, so it won't move the file between different drives.
