fix(tools): look for a NUL byte before the grep tool's built-in search reads a whole file

The built-in search of `grep` still skips a file with a NUL byte in its first 512 bytes, and now looks before it reads the rest. A file that reports a size of 0 but holds far more, such as `/proc/self/pagemap`, used to be read whole and could exhaust memory. `read`, `edit` and the built-in search now also apply their size limit to the bytes a file yields, not only to the size it reports.
