import MarkdownIt from 'markdown-it';

const markdown = new MarkdownIt({ html: false });
function inlineText(tokens = []) {
  return tokens.map(token => {
    if (token.type === 'image') return token.content;
    if (token.children) return inlineText(token.children);
    if (['text', 'code_inline'].includes(token.type)) return token.content;
    if (['softbreak', 'hardbreak'].includes(token.type)) return ' ';
    return '';
  }).join('');
}
export function plainExcerpt(source = '', limit = 200) {
  const text = markdown.parse(String(source || ''), {}).filter(token => token.type === 'inline')
    .map(token => inlineText(token.children)).join(' ').replace(/\s+/g, ' ').trim();
  const characters = Array.from(text);
  return characters.length > limit ? characters.slice(0, Math.max(0, limit)).join('') + '...' : text;
}
