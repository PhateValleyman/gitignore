gitignore-auto

> Automatically generates a .gitignore file based on your project's language(s) with Git integration.



🔧 Features

Auto-detects project type: Go, Python, Node, Java, etc.

Supports multi-language detection (e.g. Go + Python)

Optional manual override via --lang

Detects current Git branch (e.g. main, master, dev)

Downloads .gitignore templates from GitHub's github/gitignore repo

Performs Git add, commit, and push


🚀 Installation

1. Build & install (Linux/macOS)

make build     # Compile the binary
make install   # Installs to /data/data/com.termux/files/usr/bin/gitignore-auto

2. Termux (Android)

make install PREFIX=/data/data/com.termux/files/usr

🖥️ Usage

gitignore-auto

CLI Options

-h, --help         Show help message
-v, --version      Show version and author
--lang=Go,Python   Manually specify languages (skips detection)

📦 Packaging

zip gitignore-auto.zip main.go Makefile README.md

👤 Author

PhateValleyman

📧 Jonas.Ned@outlook.com



---

License: MIT

