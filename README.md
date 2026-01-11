**Vimingo**

    A lightweight, modal text editor written in Go using the Bubble Tea framework.

Vimingo brings the core Vim experience to a custom-built Go TUI. It supports standard modes (Normal, Insert, Visual, Command), syntax highlighting basics, and file management, making it a great lightweight alternative or a learning resource for building complex TUIs in Go.
✨ Features

    Modal Editing: True Vim-like experience with Normal, Insert, Visual, and Command modes.

    Motions: Navigation using standard HJKL, word jumps, and line anchors.

    Visual Mode: Highlight text to yank (copy), cut, or change.

    File Operations: robust saving and quitting workflows via Command mode.

    Status Line: Real-time feedback on current mode, filename, and cursor coordinates.

    Performance: Efficient scrolling for large files.

📦 Installation
Prerequisites

    Go 1.20+

Build from Source

    Clone the repository:
    Bash

git clone https://github.com/impossibleclone/vimingo.git
cd vimingo

Install dependencies:
Bash

go mod tidy

Build and Install:
Bash

    # Using Make (if available)
    sudo make install

    # OR manually build
    go build -o vmg main.go

🚀 Usage

Run the editor with or without a filename:
Bash

    vmg
    # OR
    vmg filename.txt

⌨️ Keybindings
🟢 Normal Mode (Navigation)

Use these keys to navigate the file.

    Key	                                                Action
    <kbd>h</kbd> <kbd>j</kbd> <kbd>k</kbd> <kbd>l</kbd>	Left, Down, Up, Right
    <kbd>w</kbd>                                        Jump to start of next word
    <kbd>e</kbd>                                        Jump to end of current word
    <kbd>_</kbd>                                        Jump to start of line
    <kbd>$</kbd>                                        Jump to end of line
    <kbd>:</kbd>                                        Enter Command Mode
    <kbd>v</kbd>                                        Enter Visual Mode

🟡 Insert Mode triggers

    Press <kbd>Esc</kbd> to exit Insert Mode.
    Key	Action
    <kbd>i</kbd>	Insert before cursor
    <kbd>a</kbd>	Insert after cursor
    <kbd>A</kbd>	Insert at end of line
    <kbd>o</kbd>	Open new line below
    <kbd>O</kbd>	Open new line above

🟣 Visual Mode

Press <kbd>v</kbd> in Normal mode to enter. Use navigation keys to highlight text.

    Key	Action
    <kbd>y</kbd>	Yank (Copy) highlighted text
    <kbd>x</kbd>	Cut highlighted text
    <kbd>c</kbd>	Change (Cut and enter Insert mode)
    <kbd>p</kbd>	Paste (works in Normal mode too)

🔴 Command Mode

Press <kbd>:</kbd> in Normal mode to enter.

    Command	Action
    :w	Save file
    :w <name>	Save as filename
    :q	Quit (fails if changes are unsaved)
    :wq	Save and Quit
    :10	Jump to line 10 (replace with any number)

🤝 Contributing

Contributions are welcome! Whether it's fixing bugs, adding new motions, or improving the TUI rendering.

    Fork the Project

    Create your Feature Branch (git checkout -b feature/AmazingFeature)

    Commit your Changes (git commit -m 'Add some AmazingFeature')

    Push to the Branch (git push origin feature/AmazingFeature)

    Open a Pull Request
