# Stackctl – Project Initializr

<img src="./public/spring.png" style="border-radius: 16px; width:80"> &nbsp;&nbsp;
<img src="./public/quarkus.png" style="border-radius: 16px; width:80"> &nbsp;&nbsp;
<img src="./public/micronaut.png" style="border-radius: 16px; width:80">


**Stackctl** is a fast, interactive Terminal User Interface (TUI) CLI for creating Java application projects from framework starters.

Select your framework, choose dependencies and project options, and generate a ready to run project directly from your terminal.

Built for a keyboard first workflow, Stackctl removes the need to open a browser just to bootstrap a new application.

Powered by Go, Bubble Tea, and official framework APIs.

**No browser. No boilerplate. Just ship. ⚡**

---

## ✨ Features

🚀 Project Initialization
  - Generate new projects using official APIs
  - Supported stacks:
    Spring Boot  
    Quarkus   
    Micronaut

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

### Install directly
## Installation

### Linux and macOS

```bash
curl -fsSL https://raw.githubusercontent.com/subrotokumar/stackctl/main/scripts/install.sh | sh
```

### Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/subrotokumar/stackctl/main/scripts/install.ps1 | iex
```

### Windows CMD

```cmd
curl -fsSL https://raw.githubusercontent.com/subrotokumar/stackctl/main/scripts/install.cmd -o install.cmd && install.cmd
```


### Install with Go

```cmd
go install github.com/subrotokumar/stackctl@latest
```

### Manual Build Installation

You can also clone the repository and run the installation script directly:

```bash
git clone https://github.com/subrotokumar/stackctl.git
cd stackctl
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