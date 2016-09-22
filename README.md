# lsel
Line selector command like pager (ex. less) written in go.

This program use some code of https://github.com/ktat/go-pager .

## Summary

1. This program take input from pipe, then show input text like pager.
2. User choose the result using hjkl and ENTER.
3. This program output result line from beginning of line to ":", that is, file path.

The shell script loated at bin/sel is normal usecase (cvim is my custom vim name, please change it as you want).

User use this program like:

$ grep SomeKeayWord *.go | sel

And you can open the result by cvim -- remote.
