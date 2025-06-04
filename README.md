Go Text Editor

A simple terminal-based text editor written in Go using the tcell package.
Features

Insert and command modes (similar to Vim)

Basic editing: typing, deleting, new lines

Command mode supports saving (:w), quitting (:q), and saving & quitting (:wq)

Scrolling support for large files

Cursor navigation with h, j, k, l keys in Normal mode

Installation

Clone the repo:

    git clone https://github.com/impossibleclone/uys.git
    cd go-text-editor

Build the project:

    go build -o uys

Run the editor:

    ./uys <filename>

Usage

Start the editor with a file name (existing or new).

Press i to enter Insert mode and start typing.

Use arrow keys or h, j, k, l to move cursor.

Press Esc to switch back to Normal mode.

Press : to enter Command mode:

:w to save

:q to quit

:wq to save and quit

Dependencies

    Go 1.20+

    tcell

Contributing

Contributions are welcome! Feel free to open issues or pull requests.
