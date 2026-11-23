import api from '../api/request'

export async function compressImageToWebP(file) {
  if (file.type === 'image/gif') return file
  return new Promise(resolve => {
    const img = new Image()
    const url = URL.createObjectURL(file)
    img.onload = () => {
      URL.revokeObjectURL(url)
      let { width, height } = img
      const ratio = Math.min(1, 2560 / Math.max(width, height))
      const canvas = document.createElement('canvas')
      canvas.width = Math.max(1, Math.round(width * ratio))
      canvas.height = Math.max(1, Math.round(height * ratio))
      const context = canvas.getContext('2d')
      if (!context) return resolve(file)
      context.drawImage(img, 0, 0, canvas.width, canvas.height)
      canvas.toBlob(blob => {
        resolve(blob && blob.size < file.size
          ? new File([blob], file.name.replace(/\.[^.]+$/, '') + '.webp', { type: 'image/webp' })
          : file)
      }, 'image/webp', 0.85)
    }
    img.onerror = () => { URL.revokeObjectURL(url); resolve(file) }
    img.src = url
  })
}

export async function uploadImageFile(file) {
  const types = { png: 'image/png', jpg: 'image/jpeg', jpeg: 'image/jpeg', gif: 'image/gif', webp: 'image/webp' }
  const type = file.type || types[file.name.split('.').pop().toLowerCase()]
  if (!Object.values(types).includes(type)) throw new Error('仅支持 PNG、JPEG、GIF、WebP 图片')
  if (!file.type) file = new File([file], file.name, { type })
  const processed = await compressImageToWebP(file)
  if (processed.size > 10 * 1024 * 1024) throw new Error('图片压缩后仍超过 10MB，请缩小后重试')
  const form = new FormData()
  form.append('file', processed)
  const result = await api.post('/upload', form)
  return result.url
}
