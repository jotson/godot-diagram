# Overview

A program that makes mermaid class diagrams from Godot source code.

This is a prototype I wrote over the weekend and probably won't work or do what you think it should do.

# Requirements

* Scenes are `.tscn` files
* Scene written in GDScript

# Usage

```bash
cd <game folder or game subfolder>
go run diagram.go
```

A file `output.mmd` will be generated in the current directory. It is a Mermaid JS diagram file that you can paste here: https://mermaid-js.github.io/mermaid-live-editor/