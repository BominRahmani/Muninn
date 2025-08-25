import './style.css';
import './app.css';
import {SaveNote} from '../wailsjs/go/main/App';
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
    if (type.startsWith('image/')) return '🖼️';
    if (type.includes('pdf')) return '📄';
    if (type.includes('text')) return '📝';
    return '📎';
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
            ID: crypto.randomUUID(), 
            Text: text,
            Attachments: attachmentInfo,
            Timestamp: new Date().toISOString() 
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

  // Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        console.log('Closing capture popup...');
        Hide(); // Hide the window when Escape is pressed
    }
    if (e.ctrlKey && e.key === 'Enter') {
        e.preventDefault();
        saveNote();
    }
});

// Make functions globally available for onclick handlers
window.triggerFileSelect = triggerFileSelect;
window.removeAttachment = removeAttachment;

// Focus text input on load
setTimeout(() => textInput.focus(), 100);
