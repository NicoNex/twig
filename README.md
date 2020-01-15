# twig
A faster version of the tree util with portability in mind.

## Usage
Usage of twig:
  -a	Prints all files including the hidden ones.
  -c	Set to false to disable colours. (default true)
  -d	Prints only the directories.

## Example:
```
$ twig foo/
foo/
├── readme.md
├── example.go
└── bar
    ├── script.sh
    ├── sas
    │   ├── sprite.jpg
    │   └── photo.png
    └── executable

2 directories, 6 files
```
