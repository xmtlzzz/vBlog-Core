import { config } from 'md-editor-v3'
import { HighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { tags } from '@lezer/highlight'

// Keep CodeMirror's behavior and theme extensions; replace only syntax colors.
const highlight = HighlightStyle.define([
  { tag: [tags.keyword, tags.bool, tags.null], color: 'var(--syntax-keyword) !important' },
  { tag: [tags.string, tags.regexp], color: 'var(--syntax-string) !important' },
  { tag: [tags.operator, tags.punctuation], color: 'var(--syntax-operator) !important' },
  { tag: [tags.number, tags.link, tags.url, tags.heading, tags.meta, tags.labelName], color: 'var(--accent) !important' },
  { tag: tags.comment, color: 'var(--muted) !important' }
])
config({ codeMirrorExtensions: extensions => [...extensions, { type: 'vblog-highlight', extension: syntaxHighlighting(highlight) }] })
