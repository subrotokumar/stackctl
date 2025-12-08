# SpringX – Spring Initializr TUI CLI 🚀

A fast, interactive **Terminal User Interface (TUI)** CLI tool for generating new **Spring Boot** projects using the official **Spring Initializr** API — powered by **Go** and the beautiful **Bubble Tea** framework.
No browser required. Just pick your project options from the terminal and bootstrap instantly! ⚡

---

## ✨ Features

* 🌀 Interactive TUI experience using Bubble Tea
* 📦 Select:

  * Spring Boot version
  * Project type (Maven/Gradle)
  * Java version
  * Dependencies (multi-select with search)
* 📁 Generates ready-to-run Spring Boot project zip
* 🛠 Download + auto-extract to project folder
* 🔗 Uses the official Spring Initializr metadata API
* 🧩 Keyboard-first navigation

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
git clone https://github.com/subrotokumar/springx.git
cd springx
go build -o springx
```

Install to PATH:

```bash
go install
```

---

## 🚀 Usage

Just run:

```bash
springx
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