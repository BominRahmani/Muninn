#!/usr/bin/osascript

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Quick Notes
# @raycast.mode silent
# @raycast.packageName Notes
#
# Optional parameters:
# @raycast.icon 📝
#
# Documentation:
# @raycast.description Instantly open vim in terminal for quick note taking
# @raycast.author Your Name
# @raycast.authorURL https://github.com/yourusername


set notesDir to (POSIX path of (path to home folder)) & "quick_notes"
set currentDate to do shell script "date +%Y%m%d_%H%M%S"
set fileName to "note_" & currentDate & ".md"
set filePath to quoted form of (notesDir & "/" & fileName)

do shell script "mkdir -p " & quoted form of notesDir

do shell script "/Applications/kitty.app/Contents/MacOS/kitty sh -c 'cd " & quoted form of notesDir & " && nvim " & quoted form of fileName & "'"


return ""
