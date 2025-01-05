# Dep Recon

Dep Recon is a tool to help quickly understand the dependencies of a project, and identify dependencies that are used for security sensitive actions (such as authorization, authentication, etc.). It is designed to make code audits quicker, as you can focus only on security sensitive dependencies, and quickly understand any DSLs and how they are used.

It works by parsing the manifest file, then downloading READMEs for each dependency from the corresponding package manager. It then parses the READMEs to keyword match against certain keywords (see `keywords.json`).

## Usage

```bash
./dep-recon -scan <path to manifest file> [-keywords <path to keywords file> -cache <path to cache directory>]
```
