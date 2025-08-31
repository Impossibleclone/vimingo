**vimingo**

A simple terminal-based text editor written in Go using the tcell package.
Features implemented till now:

Insert and command modes (similar to Vim)

Basic editing: typing, deleting, new lines

Command mode supports saving (:w), quitting (:q), and saving & quitting (:wq)

Scrolling support for large files

Motion keys:
    
    h - Left
    j - Down
    k - Up
    l - Right
    w - Next word
    e - End of Word
    _ - Start of Line
    $ - End of Line

A functional Status Line with Current Mode, File Name, Cursor Coordinates

Installation

Clone the repo:

    git clone https://github.com/impossibleclone/vimingo.git
    cd vimingo

Build the project:

    make build

Run the editor:

    ./vmg
    ./vmg <filename>

Usage

Start the editor with a file name (existing or new).

Press i to enter Insert mode and start typing.
other mappings are also available like:

    a to insert after cursor.

    A to insert at the end of line.

    o to insert at a new line below.

    O to insert at a new line above.

Use arrow keys or h, j, k, l to move cursor.

Press Esc to switch back to Normal mode.

Press v to switch to Visual mode:

In Visual mode the motions can be used to highlight the text.

    y to yank

    x to cut

    c to change

    p to paste

Press : to enter Command mode:

    :w filename.ext

    :w to save

    :q to quit

:wq to save and quit

Dependencies

    Go 1.20+

    tcell

Contributions are welcome! Feel free to open issues or pull requests.
