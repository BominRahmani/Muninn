# Project: Muninn

> A lightweight "second brain" capture client — quick notes, screenshots, and
> ideas at your fingertips.

---

## Vision

Muninn is the **capture layer** of your personal knowledge system. It’s designed
to be **frictionless**: with a single hotkey, you can jot down thoughts, paste
clipboard content, or attach screenshots. Muninn saves everything locally and
syncs with **Huginn**, the backend server, for processing and retrieval.

---

## Features (Planned)

- **Quick Capture**
  - Hotkey-triggered popup for instant notes [x]
  - Clipboard text/image capture [x]
  - Screenshot attachments [x]

- **Storage**
  - Local storage in SQLite or JSON
  - Metadata: timestamp, tags, source, attachments

- **Sync**
  - Local-first (works offline)
  - Syncs with Huginn server for processing, search, and LLM features

---

## Tech Stack

- **Language/UI**: Go + Wails
- **Storage**: SQLite or JSON (local)
- **Attachments**: Local folder with Markdown references
- **Sync**: REST API → Huginn

---

## Roadmap

### MVP
- [x] Hotkey → popup → save note (text only)
- [x] Store notes in JSON/SQLite with timestamp

### Phase 2
- [x] Clipboard text/image capture
- [x] Screenshot capture
- [x] Set up CRON job to send thoughts

### Phase 3
- [ ] Sync with Huginn server
- [x] Lightweight UI to fetch thoughts and fuzzy search on client side

---

## Relationship to Huginn

Muninn = **capture client**  
Huginn = **processing + intelligence server**

Muninn catches the notes. Huginn makes sense of them.
