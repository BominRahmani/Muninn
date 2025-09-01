import './style.css';
import './app.css';
import {SaveNote, SearchNotes, LaunchNvimWithNote} from '../wailsjs/go/main/App';
import {Hide} from '../wailsjs/runtime/runtime';

// Capture functionality
const popup = document.getElementById('capturePopup');
const textInput = document.getElementById('textInput');
const attachmentsSection = document.getElementById('attachmentsSection');
const fileInput = document.getElementById('fileInput');
const pastedIndicator = document.getElementById('pastedIndicator');

let attachments = [];

// Auto-resize textarea
textInput.addEventListener('input', function() {
    this.style.height = 'auto';
    this.style.height = Math.max(120, this.scrollHeight) + 'px';
});

// Drag and drop handlers
popup.addEventListener('dragover', (e) => {
    e.preventDefault();
    popup.classList.add('drag-over');
});

popup.addEventListener('dragleave', (e) => {
    if (!popup.contains(e.relatedTarget)) {
        popup.classList.remove('drag-over');
    }
});

popup.addEventListener('drop', (e) => {
    e.preventDefault();
    popup.classList.remove('drag-over');
    
    const files = Array.from(e.dataTransfer.files);
    addAttachments(files);
});

// Paste handler
document.addEventListener('paste', (e) => {
    const items = e.clipboardData.items;
    const files = [];
    
    for (let item of items) {
        if (item.kind === 'file') {
            files.push(item.getAsFile());
        }
    }
    
    if (files.length > 0) {
        addAttachments(files);
        showPastedIndicator();
    }
});

// File input handler
fileInput.addEventListener('change', (e) => {
    const files = Array.from(e.target.files);
    addAttachments(files);
    fileInput.value = ''; // Reset input
});

function triggerFileSelect() {
    fileInput.click();
}

function addAttachments(files) {
    files.forEach(file => {
        const attachment = {
            id: Date.now() + Math.random(),
            name: file.name,
            file: file,
            type: file.type
        };
        attachments.push(attachment);
    });
    renderAttachments();
}

function removeAttachment(id) {
    attachments = attachments.filter(att => att.id !== id);
    renderAttachments();
}

function renderAttachments() {
    if (attachments.length === 0) {
        attachmentsSection.classList.remove('has-attachments');
        attachmentsSection.innerHTML = '';
        return;
    }

    attachmentsSection.classList.add('has-attachments');
    attachmentsSection.innerHTML = attachments.map(att => `
        <div class="attachment-item">
            <div class="attachment-icon">${getFileIcon(att.type)}</div>
            <div class="attachment-name">${att.name}</div>
            <button class="attachment-remove" onclick="removeAttachment(${att.id})">×</button>
        </div>
    `).join('');
}

function getFileIcon(type) {
    if (type.startsWith('image/')) return 'IMG';
    if (type.includes('pdf')) return 'PDF';
    if (type.includes('text')) return 'TXT';
    return 'FILE';
}

function showPastedIndicator() {
    pastedIndicator.classList.add('show');
    setTimeout(() => {
        pastedIndicator.classList.remove('show');
    }, 2000);
}

async function saveNote() {
    const text = textInput.value.trim();
    if (!text && attachments.length === 0) {
        return;
    }

    const capturePopup = document.getElementById('capturePopup');
    capturePopup.classList.add('screen-vibrate');
    setTimeout(() => capturePopup.classList.remove('screen-vibrate'), 300);

    try {
        const attachmentInfo = await Promise.all(
            attachments.map(async (a, i) => {
                let filePath = '';
                let data = null;
                
                // Convert file to bytes for backend
                const file = a.file;
                if (file.path) {
                    filePath = file.path; // For files with system paths
                } else if (file instanceof File || file instanceof Blob) {
                    const arrayBuffer = await file.arrayBuffer();
                    data = Array.from(new Uint8Array(arrayBuffer)); // Convert to byte array
                }

                return {
                    id: "", 
                    fileName: a.name || `attachment_${i}`,
                    fileType: a.type || "application/octet-stream",
                    filePath, 
                    data // Send bytes to backend
                };
            })
        );

        const thought = {
            id: crypto.randomUUID(), 
            text: text,
            attachments: attachmentInfo,
            timestamp: new Date().toISOString() 
        };

        await SaveNote(thought);
        
        // Reset UI
        textInput.value = '';
        textInput.style.height = '120px';
        attachments = [];
        renderAttachments();

        console.log('Note saved successfully');
    } catch (err) {
        console.error('Error saving note:', err);
    }
}

const input = document.getElementById("search-input");
const resultsList = document.getElementById("search-results");
const overlay = document.getElementById("search-overlay");

let debounceTimer;
let selectedIndex = -1;
let searchResults = [];
let justShowedPreview = false;

input.addEventListener("input", () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    performSearch(input.value);
  }, 150);
});

// Keyboard navigation
input.addEventListener("keydown", (e) => {
  if (e.key === "ArrowDown") {
    e.preventDefault();
    navigateResults(1);
  } else if (e.key === "ArrowUp") {
    e.preventDefault();
    navigateResults(-1);
  } else if (e.key === "Enter") {
    e.preventDefault();
    selectCurrentResult();
  } else if (e.key === "Escape") {
    hideAllOverlays();
  }
});

function navigateResults(direction) {
  const items = resultsList.querySelectorAll('.search-result-item');
  if (items.length === 0) return;

  // Remove current selection
  if (selectedIndex >= 0 && selectedIndex < items.length) {
    items[selectedIndex].classList.remove('selected');
  }

  // Update selection
  selectedIndex += direction;
  if (selectedIndex < 0) selectedIndex = items.length - 1;
  if (selectedIndex >= items.length) selectedIndex = 0;

  // Add selection to new item
  if (selectedIndex >= 0 && selectedIndex < items.length) {
    items[selectedIndex].classList.add('selected');
    items[selectedIndex].scrollIntoView({ block: 'nearest' });
  }
}

function selectCurrentResult() {
  const items = resultsList.querySelectorAll('.search-result-item');
  if (selectedIndex >= 0 && selectedIndex < items.length) {
    items[selectedIndex].click();
  }
}

async function performSearch(query) {
  if (!query.trim()) {
    resultsList.innerHTML = "";
    searchResults = [];
    selectedIndex = -1;
    return;
  }

  // Show loading state
  resultsList.innerHTML = `
    <li class="loading-state">
      <div class="loading-icon">⏳</div>
      <div class="loading-text">Searching...</div>
    </li>
  `;

  try {
    const results = await SearchNotes(query);
    searchResults = results;
    
    resultsList.innerHTML = "";
    selectedIndex = -1;
    
    if (results.length === 0) {
      const noResults = document.createElement("li");
      noResults.className = "no-results";
      noResults.innerHTML = `
        <div class="no-results-icon">SEARCH</div>
        <div class="no-results-text">No results found</div>
      `;
      resultsList.appendChild(noResults);
      return;
    }

    results.forEach((note, index) => {
      const li = document.createElement("li");
      li.className = "search-result-item";
      li.innerHTML = `
        <div class="result-icon">📝</div>
        <div class="result-content">
          <div class="result-text">${highlightQuery(note.content, query)}</div>
          <div class="result-meta">${formatTimestamp(note.timestamp)} • ${note.content.length > 200 ? 'Long note' : 'Short note'}</div>
        </div>
      `;
      li.onclick = () => {
        console.log("Selected note:", note.id);
        showNotePreview(note);
      };
      resultsList.appendChild(li);
    });
    
    // Update footer with results count
    updateSearchFooter(results.length, query);
  } catch (err) {
    console.error("Search failed", err);
    resultsList.innerHTML = `
      <li class="error-message">
        <div class="error-icon">⚠️</div>
        <div class="error-text">Search failed</div>
      </li>
    `;
  }
}

function highlightQuery(text, query) {
  if (!query) return text;
  
  const regex = new RegExp(`(${query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
  return text.replace(regex, '<mark>$1</mark>');
}

function formatTimestamp(timestamp) {
  if (!timestamp) return '';
  
  const date = new Date(timestamp);
  const now = new Date();
  const diffInHours = (now - date) / (1000 * 60 * 60);
  
  if (diffInHours < 1) {
    return 'Just now';
  } else if (diffInHours < 24) {
    return `${Math.floor(diffInHours)}h ago`;
  } else {
    return date.toLocaleDateString();
  }
}

function showCopyFeedback() {
  const feedback = document.createElement('div');
  feedback.className = 'copy-feedback';
  feedback.textContent = 'Copied to clipboard!';
  document.body.appendChild(feedback);
  
  setTimeout(() => {
    feedback.remove();
  }, 2000);
}

function updateSearchFooter(resultCount, query) {
  const footer = document.querySelector('.search-shortcuts');
  if (footer) {
    if (resultCount > 0) {
      footer.innerHTML = `
        <span class="results-count">${resultCount} result${resultCount === 1 ? '' : 's'}</span> • 
        <span class="shortcut">⌘K</span> to search • 
        <span class="shortcut">Click</span> to preview • 
        <span class="shortcut">Backspace</span> to close
      `;
    } else {
      footer.innerHTML = `
        <span class="shortcut">⌘K</span> to search • 
        <span class="shortcut">Click</span> to preview • 
        <span class="shortcut">Backspace</span> to close
      `;
    }
  }
}

function showCaptureOverlay() {
  document.getElementById("search-overlay").classList.add("hidden");
  document.getElementById("capture-overlay").classList.remove("hidden");
  document.getElementById("textInput").focus();
}

function showSearchOverlay() {
  document.getElementById("capture-overlay").classList.add("hidden");
  document.getElementById("search-overlay").classList.remove("hidden");
  document.getElementById("search-input").focus();
}

function hideAllOverlays() {
  document.getElementById("capture-overlay").classList.add("hidden");
  document.getElementById("search-overlay").classList.add("hidden");
  document.getElementById("note-preview-overlay").classList.add("hidden");
}

function showNotePreview(note) {
  // Hide search overlay
  document.getElementById("search-overlay").classList.add("hidden");
  
  // Populate note preview content
  const noteText = document.getElementById("notePreviewText");
  const noteAttachments = document.getElementById("notePreviewAttachments");
  const noteMeta = document.getElementById("notePreviewMeta");
  
  // Set note text
  noteText.textContent = note.text || note.content || '';
  
  // Set note metadata
  const timestamp = note.timestamp ? formatTimestamp(note.timestamp) : '';
  const wordCount = (note.text || note.content || '').split(/\s+/).filter(word => word.length > 0).length;
  noteMeta.innerHTML = `
    ${timestamp} • ${wordCount} word${wordCount === 1 ? '' : 's'} • Note ID: ${note.id} • 
    <span class="shortcut">Enter</span> to edit in nvim • 
    <span class="shortcut">⌘C</span> to copy • 
    <span class="shortcut">Backspace</span> to return to search
  `;
  
  // Handle attachments if they exist
  if (note.attachments && note.attachments.length > 0) {
    noteAttachments.innerHTML = note.attachments.map(attachment => `
      <div class="note-attachment-item">
        <div class="note-attachment-icon">${getFileIcon(attachment.fileType || attachment.type || '')}</div>
        <div class="note-attachment-name">${attachment.fileName || attachment.name || 'Unknown file'}</div>
      </div>
    `).join('');
  } else {
    noteAttachments.innerHTML = '';
  }
  
  document.getElementById("note-preview-overlay").classList.remove("hidden");
  justShowedPreview = true;
  setTimeout(() => {
    justShowedPreview = false;
  }, 100);
}

function hideNotePreview() {
  document.getElementById("note-preview-overlay").classList.add("hidden");
}

function copyNoteToClipboard() {
  const noteText = document.getElementById("notePreviewText");
  const textToCopy = noteText.textContent;
  
  navigator.clipboard.writeText(textToCopy).then(() => {
    showCopyFeedback();
  }).catch(err => {
    console.error('Failed to copy note:', err);
  });
}

async function launchNvimWithCurrentNote() {
  const noteText = document.getElementById("notePreviewText");
  const textToCopy = noteText.textContent;
  
  try {
    await LaunchNvimWithNote(textToCopy);
    console.log('Launched nvim with note content');
    
    // Return to search screen after launching nvim
    hideNotePreview();
    showSearchOverlay();
  } catch (err) {
    console.error('Failed to launch nvim:', err);
    // Show error feedback to user
    const feedback = document.createElement('div');
    feedback.className = 'copy-feedback error';
    feedback.textContent = 'Failed to launch nvim';
    document.body.appendChild(feedback);
    
    setTimeout(() => {
      feedback.remove();
    }, 3000);
  }
}

// Listen for Go event to focus search
window.runtime.EventsOn("focusSearch", () => {
  showSearchOverlay();
});

window.runtime.EventsOn("focusCapture", () => {
  showCaptureOverlay();
});


// Global keyboard shortcuts
document.addEventListener('keydown', (e) => {
    if (e.key === 'Backspace') {
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA' || e.target.isContentEditable) {
            return;
        }
        
        const notePreviewOverlay = document.getElementById('note-preview-overlay');
        const searchOverlay = document.getElementById('search-overlay');
        const captureOverlay = document.getElementById('capture-overlay');
        
        if (!notePreviewOverlay.classList.contains('hidden')) {
            e.preventDefault();
            hideNotePreview();
            showSearchOverlay(); // Return to search screen when pressing backspace from preview
        } else if (!searchOverlay.classList.contains('hidden')) {
            e.preventDefault();
            hideAllOverlays();
        } else if (!captureOverlay.classList.contains('hidden')) {
            e.preventDefault();
            hideAllOverlays();
        } else {
            e.preventDefault();
            console.log('Closing capture popup...');
            Hide(); // Hide the window when Backspace is pressed
        }
    }
    
    // Ctrl+Enter to save note (only in capture mode)
    if (e.ctrlKey && e.key === 'Enter' && overlay.classList.contains('hidden')) {
        e.preventDefault();
        saveNote();
    }
    
    // Cmd/Ctrl+K to focus search
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        showSearchOverlay();
    }
    
    // Cmd/Ctrl+C to copy note (when in note preview)
    if ((e.metaKey || e.ctrlKey) && e.key === 'c') {
        const notePreviewOverlay = document.getElementById('note-preview-overlay');
        if (!notePreviewOverlay.classList.contains('hidden')) {
            e.preventDefault();
            copyNoteToClipboard();
        }
    }
    
    // Enter key to launch nvim (when in note preview)
    if (e.key === 'Enter') {
        const notePreviewOverlay = document.getElementById('note-preview-overlay');
        if (!notePreviewOverlay.classList.contains('hidden') && !justShowedPreview) {
            e.preventDefault();
            launchNvimWithCurrentNote();
        }
    }
});

window.triggerFileSelect = triggerFileSelect;
window.removeAttachment = removeAttachment;
window.showNotePreview = showNotePreview;
window.hideNotePreview = hideNotePreview;
window.copyNoteToClipboard = copyNoteToClipboard;
window.launchNvimWithCurrentNote = launchNvimWithCurrentNote;

document.addEventListener('DOMContentLoaded', () => {
    const editNvimBtn = document.getElementById('editNvimBtn');
    const copyNoteBtn = document.getElementById('copyNoteBtn');
    const closeNoteBtn = document.getElementById('closeNoteBtn');
    
    if (editNvimBtn) {
        editNvimBtn.addEventListener('click', launchNvimWithCurrentNote);
    }
    
    if (copyNoteBtn) {
        copyNoteBtn.addEventListener('click', copyNoteToClipboard);
    }
    
    if (closeNoteBtn) {
        closeNoteBtn.addEventListener('click', () => {
            hideNotePreview();
            showSearchOverlay(); // Return to search screen when closing via button
        });
    }
});

setTimeout(() => textInput.focus(), 100);
