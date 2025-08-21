# Project: Muninn 

> A lightweight "second brain" for quick capture, organization, and retrieval of
> notes, screenshots, and ideas — powered by NLP and semantic search.

---

##  Vision

Muninn is a personal knowledge capture system designed to be **frictionless** for
note-taking and **powerful** for retrieval. With a single hotkey, you can jot
down thoughts, paste clipboard content, or attach screenshots. Over time, Muninn
aggregates these into a structured knowledge base, enriched with
auto-categorization, semantic search, and intelligent linking.


---

## Features (Planned)

- **Quick Capture**
  - Hotkey-triggered popup for instant notes
  - Clipboard text/image capture
  - Screenshot attachments

- **Storage**
  - Notes stored as structured JSON (or SQLite)
  - Metadata: timestamp, tags, source, attachments

- **Processing**
  - Auto-categorization (rule-based + NLP)
  - Semantic similarity linking (Markdown-style `[[links]]`)
  - Fuzzy + semantic search across all notes

- **Retrieval**
  - Search bar with fuzzy + embedding-based search
  - Export to Markdown with links
  - Optional weekly digest / resurfacing of old notes

- **Sync (Future)**
  - Git-based sync (simple)
  - Optional FastAPI server for multi-device sync

## LLM Capabilities (Planned)

Large Language Models (LLMs) can supercharge Seneca by adding intelligence on
top of the raw notes. While the core system works offline with embeddings and
search, LLMs provide **categorization, summarization, and natural-language
interaction**.

### Where LLMs Fit

- **Smart Categorization & Tagging**
  - Suggest categories or tags for new notes
  - Detect todos, questions, or references automatically

- **Summarization & Digest**
  - Generate daily/weekly summaries of new notes
  - Highlight recurring themes and patterns
  - Provide a "reflection mode" for reviewing knowledge

- **Semantic Q&A (Chat with Your Notes)**
  - Ask natural questions like:
    - "What were my last 3 ideas about Project X?"
    - "Show me all notes related to debugging Python scripts."
  - Uses a RAG (Retrieval-Augmented Generation) pipeline:
    1. Embed query → retrieve top notes
    2. Pass notes + query to LLM
    3. LLM generates a natural-language answer with references

- **Note Linking & Relationship Discovery**
  - Explain why two notes are related
  - Insert contextual Markdown links:
    ```markdown
    [[Note123]] is related to [[Note456]] (both discuss debugging strategies).
    ```

- **Knowledge Expansion**
  - Expand short notes into more detailed ideas
  - Example: "Idea: automate screenshot capture for bug reports" → expanded into
    implementation suggestions and libraries

---

## Tech Stack (Proposed)

- **Language**: Python/Go
- **UI**: Tauri
- **Storage**: SQLite or JSON -> Saved on PC that I need to convert to server
- **NLP**: `sentence-transformers`, `scikit-learn`
- **LLM**: OpenAI API
- **Search**: FAISS (semantic), `rapidfuzz` (fuzzy)
- **Attachments**: Local folder with Markdown references
- **Sync**: Git or FastAPI server

---

## TODO Roadmap

### MVP
- [ ] Hotkey → popup → save note (text only)
- [ ] Store notes in JSON/SQLite with timestamp
- [ ] Basic fuzzy search

### Phase 2
- [ ] Add clipboard text/image capture
- [ ] Add screenshot capture
- [ ] Add tagging via `#hashtags`

### Phase 3
- [ ] Generate embeddings for notes
- [ ] Implement semantic search (FAISS or similar)
- [ ] Auto-categorization (rule-based + NLP)
- [ ] Auto-linking of related notes (`[[NoteID]]`)

### Phase 4
- [ ] Export notes to Markdown with links
- [ ] Weekly digest / resurfacing old notes
- [ ] Git-based sync

### Phase 5
- [ ] FastAPI server for multi-device sync
- [ ] Advanced clustering / summarization
- [ ] Optional graph visualization


### Phase 6
- [ ] Add LLM capabilities

---


