// Shared component detection script and srcdoc builder for iframe-based component preview.
// Used by: admin/Custom.vue, blog/Modules.vue, shared/CustomWidgets.vue

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
  'if(__fns.length){try{__fns[0].fn(c)}catch(ex){document.body.innerHTML="<pre style=\\"color:red;padding:12px\\">Error calling "+__fns[0].name+": "+ex.message+"</pre>"}}' +
  'else{document.body.innerHTML="<pre style=\\"color:orange;padding:12px\\">No component function found. Code may not define a function accepting a container argument.</pre>"}' +
  '}catch(e){document.body.innerHTML="<pre style=\\"color:red;padding:12px\\">"+e.message+"\\n"+(e.stack||"")+"</pre>"}'

const BASE_STYLES =
  'body{margin:0;font-family:system-ui,sans-serif}' +
  ':root{--bg:#fafafa;--surface:#fff;--fg:#171717;--muted:#737373;--border:#e5e5e5;--accent:#2563eb;--success:#16a34a;--font-sans:system-ui,sans-serif;--radius:8px}'

export function buildSrcdoc(code, { padding = '12px', data = null } = {}) {
  if (!code) return ''
  var dataScript = data ? '<script>window.__BLOG_DATA__=' + JSON.stringify(data) + '<\/script>' : ''
  return (
    '<!DOCTYPE html><html><head><style>' +
    BASE_STYLES +
    'body{padding:' + padding + '}' +
    '</style></head><body><div id="root"></div>' +
    dataScript +
    '<script>' + code + '<\/script>' +
    '<script>' + DETECT_SCRIPT + '<\/script>' +
    '</body></html>'
  )
}
