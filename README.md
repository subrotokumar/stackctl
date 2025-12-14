# Stackctl – Java Project Initializr TUI CLI 🚀

Stackctl is a fast, interactive Terminal User Interface (TUI) CLI for managing the entire lifecycle of Java-based projects.

It goes beyond project generation by offering a stack of developer utilities — project initialization, testing, scripts, and automation — all from a clean, keyboard-first terminal experience.

Powered by Go, Bubble Tea, and official framework APIs.

No browser. No boilerplate. Just ship. ⚡

---

## ✨ Features

🚀 Project Initialization
  - Generate new projects using official APIs
  - Supported stacks:  
    Spring Boot  
    Quarkus   
    Micronaut (planned)  

Configure:
  - Build tool (Maven / Gradle)
  - Java version
  - Framework version 

Dependencies (multi-select with search)

---

## 📸 Preview

![](./public/demo.gif)


---

## 🧰 Tech Stack

| Component             | Purpose                       |
| --------------------- | ----------------------------- |
| Go                    | Core language                 |
| Bubble Tea            | TUI UI framework              |
| Lipgloss              | Styling                       |
| Bubbles               | UI widgets                    |
| Spring Initializr API | Metadata + project generation |

---

## 📦 Installation

### From Source

Make sure you have Go 1.21+ installed.

```bash
git clone https://github.com/subrotokumar/stackctl.git
cd stackctl
go build -o stackctl
```

Install to PATH:

```bash
go install
```

---

## 🚀 Usage

Just run:

```bash
stackctl
```

Follow the interactive terminal UI to configure your project.
Once done, your Spring Boot project will be created and extracted automatically.

---

## ⌨️ Controls

| Action            | Key        |
| ----------------- | ---------- |
| Navigate          | ↑ ↓ or j k |
| Select / Continue | Enter      |
| Go Back           | Esc        |
| Multi Select      | Space      |
| Quit              | Ctrl + C   |

---

## 🤝 Contributing

Contributions are welcome! ❤️
Please open an issue or submit a pull request.

---

## 📝 License

This project is licensed under the **MIT License**.
See the `LICENSE` file for more details.

---

## ⭐ Support

If you like this project:

* ⭐ Star the repo
* 🔁 Share with other Spring + Go developers