# Overview

A program that makes mermaid class diagrams from Godot source code.

This is a prototype I wrote over the weekend and probably won't work or do what you think it should do.

![Example Mermaid diagram](https://raw.githubusercontent.com/jotson/godot-diagram/main/example.png)

# Requirements

* Scenes are `.tscn` files
* Scene code written in GDScript

# Usage

```bash
cd <game folder or game subfolder>
go run /path/to/godot-diagram/diagram.go
```

A file `output.mmd` will be generated in the current directory. It is a Mermaid JS diagram file that you can paste here: https://mermaid-js.github.io/mermaid-live-editor/