// Shared component detection script and srcdoc builder for iframe-based component preview.
// Used by: admin/Custom.vue, blog/Modules.vue, shared/CustomWidgets.vue
import themeTokens from '../styles/tokens.css?inline'

export const DETECT_SCRIPT =
  'try{' +
  'var c=document.getElementById("root");' +
  'var __skip=/^(addEventListener|removeEventListener|dispatchEvent|constructor|toString|valueOf|hasOwnProperty|isPrototypeOf|propertyIsEnumerable|toLocaleString|getComputedStyle|getSelection|matchMedia|open|close|stop|alert|confirm|prompt|fetch|btoa|atob|setTimeout|setInterval|clearTimeout|clearInterval|requestAnimationFrame|cancelAnimationFrame|reportError|structuredClone|queueMicrotask|getScreenDetails|queryLocalFonts|showDirectoryPicker|showOpenFilePicker|showSaveFilePicker|webkitRequestAnimationFrame|webkitCancelAnimationFrame|requestVideoFrameCallback|cancelIdleCallback|requestIdleCallback|createImageBitmap|scroll|scrollTo|scrollBy|moveTo|moveBy|resizeTo|resizeBy|postMessage|focus|blur|close|stop|print|find|getSelection)$/;' +
  'var __native=/^(cancel|get|set|on|request|webkit|moz|ms|create|has|is|to|from|add|remove|dispatch|match|replace|search|split|at|concat|fill|filter|find|flat|forEach|includes|indexOf|join|keys|lastIndexOf|map|pop|push|reduce|reverse|shift|slice|some|sort|splice|unshift|values|entries|every|copyWithin|length|name|prototype|caller|arguments|apply|bind|call)$/;' +
  'var __fns=[];' +
  'for(var k in window){' +
  'if(__skip.test(k))continue;' +
  'try{' +
  'if(typeof window[k]==="function"&&window[k].length>=1){' +
  'var s=window[k].toString();' +
  'if(s.indexOf("[native ")===-1)__fns.push({name:k,fn:window[k]})' +
  '}}catch(x){}}' +
  'if(!__fns.length)for(var k in window){' +
  'if(__skip.test(k)||__native.test(k))continue;' +
  'try{' +
  'if(typeof window[k]==="function"&&window[k].length>=1){' +
  'var s=window[k].toString();' +
  'if(s.indexOf("[native ")===-1)__fns.push({name:k,fn:window[k]})' +
  '}}catch(x){}}' +
  'if(__fns.length){try{__fns[0].fn(c)}catch(ex){document.body.innerHTML="<pre style=\\"color:var(--error);padding:12px\\">Error calling "+__fns[0].name+": "+ex.message+"</pre>"}}' +
  'else{document.body.innerHTML="<pre style=\\"color:var(--warning);padding:12px\\">No component function found. Code may not define a function accepting a container argument.</pre>"}' +
  '}catch(e){document.body.innerHTML="<pre style=\\"color:var(--error);padding:12px\\">"+e.message+"\\n"+(e.stack||"")+"</pre>"}'

const BASE_STYLES = themeTokens + 'body{margin:0;font-family:var(--font-sans);color:var(--fg);background:var(--bg)}'
const THEME_SCRIPT = 'window.addEventListener("message",function(event){if(event.source!==parent||event.data?.type!=="vblog-theme")return;document.documentElement.dataset.theme=event.data.theme==="dark"?"dark":"light"})'

export function buildSrcdoc(code, { padding = '12px', data = null, theme = 'light' } = {}) {
  if (!code) return ''
  var dataScript = data ? '<script>window.__BLOG_DATA__=' + JSON.stringify(data) + '<\/script>' : ''
  return (
    '<!DOCTYPE html><html data-theme="' + (theme === 'dark' ? 'dark' : 'light') + '"><head><style>' +
    BASE_STYLES +
    'body{padding:' + padding + '}' +
    '</style></head><body><div id="root"></div>' +
    dataScript +
    '<script>' + THEME_SCRIPT + '<\/script>' +
    '<script>' + code + '<\/script>' +
    '<script>' + DETECT_SCRIPT + '<\/script>' +
    '</body></html>'
  )
}
