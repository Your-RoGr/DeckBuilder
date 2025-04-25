[![Version](https://img.shields.io/badge/Version-0.1.0-blue)](https://github.com/Your-RoGr/DeckBuilder/tree/master)
[![Latest Release](https://img.shields.io/github/v/release/Your-RoGr/DeckBuilder)](https://github.com/Your-RoGr/DeckBuilder/releases)
![License](https://img.shields.io/github/license/Your-RoGr/DeckBuilder)
![Downloads](https://img.shields.io/github/downloads/Your-RoGr/DeckBuilder/total)
[![Go Report Card](https://goreportcard.com/badge/Your-RoGr/DeckBuilder)](https://goreportcard.com/report/github.com/Your-RoGr/DeckBuilder)
![GitHub Stars](https://img.shields.io/github/stars/Your-RoGr/DeckBuilder?style=social)

# DeckBuilder

[English](README.md) | [Русский](README.ru.md)

A terminal-based application to help you create and manage decks for [Anki](https://apps.ankiweb.net/).

- [Features](#Features)
- [Usage](#Usage)
- [Dependencies](#Dependencies)
- [Installation](#Installation)
- [License](#License)

## Features

- Terminal menu interface for managing decks and word lists
- Create new decks (CSV files) and add them to your catalog
- Browse and select decks from a catalog stored locally
- Add words or word-translation pairs to your decks interactively, with duplicate checking
- View and browse the contents of decks
- Delete decks from your catalog list
- File navigation and new deck creation via an interactive file chooser
- All input and navigation happens in the terminal, with keyboard hotkeys
- Decks are saved as CSV files for easy import into Anki or further processing

## Usage

To start DeckBuilder, run the app from your terminal:

```bash
DeckBuilder
```

**Menu Navigation:**
- Use the ▲ and ▼ arrow keys to move between menu options
- Press `Enter` to select an option
- Press `Esc` to go back or exit
- In the file selection menu:
  - Press `A` to create a new deck (CSV file) in the current directory
- In deck menus:
  - Press `D` to delete a deck from the catalog

**Typical Workflow:**
- Select "Select new file" to create or choose a deck file.
- Once a deck is selected, pick a mode:
  - **Word**: Add single words
  - **Word-Translate**: Add word-translation pairs
  - **Show**: Browse the contents of the deck

Entries added are unique per deck—duplicate entries are detected and rejected.

**Deck Format:**
- Each deck is saved as a CSV file, suitable for import into Anki or as a source for further processing.
- The default columns are “Word” and “Translation”, but you can use either single-word or word-translation formats.

## Dependencies

This package uses the following libraries:
- [termbox-go](https://github.com/nsf/termbox-go) – for terminal GUI
- [mattn/go-runewidth](https://github.com/mattn/go-runewidth) – for character width handling in the terminal

All other functionality is implemented with Go standard library packages.

## Installation

To install the app, follow these steps:

1. **Clone the repository:**

```bash
git clone https://github.com/Your-RoGr/DeckBuilder.git
cd DeckBuilder
```

2. **Run the build script:**  
This script will automatically build the project and install the executable into `/usr/local/bin` (you may be prompted for your password).

```bash
./scripts/build.sh
```

> **Note:**  
> - The script requires [Go](https://golang.org/doc/install) to be installed on your system.
> - The install step uses `sudo` to copy the binary to `/usr/local/bin`.

3. **Run DeckBuilder from anywhere:**

```bash
DeckBuilder
```

If you see the terminal menu, the installation was successful!

## License

DeckBuilder is MIT-Licensed
