# SVT Text CLI

Read Swedish news from SVT Text in the terminal. Using the SVT Text API (https://www.svt.se/text-tv/api).

## Example

**Interactive mode with colors**

```sh
$ svttext --colors --interactive
```

![Example](./doc/example.gif)

## Install

Find the binary for your OS under [releases](https://github.com/ahrberg/svttext/releases). Download and run the binary in your terminal of choice.

## Usage

```sh
$ svttext --help
Usage: svttext [OPTION]... [PAGE]
Read news from SVT Text

Example: svttext --colors 100

Options:
  -colors
        colorize the output
  -interactive
        start interactive mode
        use arrow keys to navigate pages
        or enter page number to go to page
  -version
        prints svttext version
```

Make an alias if you like. Bash example, place in your .bashrc:

```bash
alias svt="svttext --colors --interactive"
```

## TODO

- [ ] Add SVT logo? The "TV" version have a logo at the start page. Maybe add this as ascii art?
- [ ] Release using homebrew tap? https://appliedgo.net/release2/
- [ ] Sign releases? https://appliedgo.net/release/
- [x] Make release, github action
- [x] Colorize output, add flag `--color`
- [x] Interactive mode, add a flag `--interactive` and don't exit the cli instead give options like enter new page, navigate with arrow keys, vi key binding etc.
