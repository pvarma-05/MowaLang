'use client';

import * as monaco from 'monaco-editor';

export function defineMowaTheme() {
  if (typeof window === 'undefined') return; // Skip on server

  monaco.editor.defineTheme('mowaTheme', {
    base: 'vs',
    inherit: true,
    rules: [
      { token: 'keyword', foreground: '#A93E39', fontStyle: 'bold' }, // Keywords like idhi, theesko
      { token: 'string', foreground: '#006400' }, // Strings
      { token: 'number', foreground: '#0000FF' }, // Numbers
      { token: 'comment', foreground: '#6A9955' }, // Comments
      { token: 'identifier', foreground: '#000000' }, // Variables
    ],
    colors: {
      'editor.background': '#FDF2EC', // Light peach background
      'editor.foreground': '#000000', // Black text
      'editorLineNumber.foreground': '#A93E39', // Line numbers in red
      'editorLineNumber.activeForeground': '#7B2D2A', // Active line number
      'editor.selectionBackground': '#A93E3933', // Selection with opacity
      'editorCursor.foreground': '#A93E39', // Cursor in red
    },
  });
}
