# fscan4go
A file scancer using golang, list(recursive) file's size and type on disk, simple but fast.

### Usage:
```
Usage:
    ./fsc4g PATH 
    ./fsc4g PATH REGEX 
Example:
    ./fsc4g /somewhere/tmp "^.*\.txt$" 
```

Recursive searching contains hidden folders and list all files if no REGEX specified.


