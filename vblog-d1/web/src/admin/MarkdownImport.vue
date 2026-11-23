<template>
  <input ref="markdownInput" type="file" accept=".md,.markdown" hidden @change="receiveFiles($event, false)" />
  <input ref="directoryInput" type="file" webkitdirectory multiple hidden @change="receiveFiles($event, true)" />
  <input ref="extraDirectoryInput" type="file" webkitdirectory multiple hidden @change="addFiles" />
  <input ref="extraFilesInput" type="file" accept=".png,.jpg,.jpeg,.gif,.webp" multiple hidden @change="addFiles" />
  <input ref="replacementInput" type="file" accept=".png,.jpg,.jpeg,.gif,.webp" hidden @change="replaceFile" />
  <el-dialog v-model="visible" title="导入 Markdown 文章" width="min(760px, calc(100vw - 32px))"
    :close-on-click-modal="!busy" :close-on-press-escape="!busy" :show-close="!busy">
    <p class="import-help">选择包含文章和图片的文件夹，或上传 Markdown 后补选图片文件夹。相对路径对应且结果唯一的图片会自动选中，上传后自动替换正文中的路径。只上传文章引用到的图片。</p>
    <label v-if="markdowns.length > 1" class="article-choice">选择文章
      <select v-model="selectedPath" :disabled="busy || reading" @change="loadSelected">
        <option value="">请选择要导入的 Markdown</option>
        <option v-for="file in markdowns" :key="filePath(file)" :value="filePath(file)">{{ filePath(file) }}</option>
      </select>
    </label>
    <p v-else-if="markdownFile" class="import-filename">{{ filePath(markdownFile) }}</p>
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />
    <p v-if="reading" role="status">正在读取文章…</p>
    <template v-else-if="markdownFile && source !== null">
      <div v-if="images.length" class="import-toolbar">
        <span>{{ images.length }} 张本地图片<span v-if="missingCount"> · {{ missingCount }} 张待补选或确认</span></span>
        <div>
          <el-button size="small" :disabled="busy" @click="extraDirectoryInput.click()">补选图片文件夹</el-button>
          <el-button size="small" :disabled="busy" @click="extraFilesInput.click()">补选图片</el-button>
        </div>
      </div>
      <div v-if="images.length" class="image-list">
        <div v-for="(image, index) in images" :key="image.source" class="image-row">
          <div class="image-source">
            <code>{{ image.source }}</code>
            <span :class="{ 'image-error': image.error }" role="status">{{ image.error || statusLabel(image) }}</span>
          </div>
          <div class="image-match">
            <select :value="image.file ? files.indexOf(image.file) : ''" :disabled="busy" :aria-label="'匹配图片：' + image.source" @change="chooseFile(image, $event)">
              <option value="">{{ image.candidates.length ? '请确认对应的图片' : '未找到图片，请补选' }}</option>
              <option v-for="file in optionsFor(image)" :key="files.indexOf(file)" :value="files.indexOf(file)">{{ filePath(file) }}</option>
            </select>
            <el-button size="small" :disabled="busy" @click="selectReplacement(index)">选择图片</el-button>
          </div>
        </div>
      </div>
      <p v-else class="import-help">未发现需要上传的本地图片，可以直接进入编辑。</p>
      <el-progress v-if="busy && images.length" :percentage="Math.round(images.filter(image => image.status === 'uploaded').length / images.length * 100)" />
    </template>
    <template #footer>
      <el-button :disabled="busy" @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="busy" :disabled="reading || !markdownFile || source === null || missingCount > 0" @click="finish">
        {{ images.length ? '上传图片并进入编辑' : '进入编辑' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, ref, shallowRef } from 'vue'
import { uploadImageFile } from '../utils/imageUpload'

const emit = defineEmits(['imported'])
const markdownInput = ref(null), directoryInput = ref(null), extraDirectoryInput = ref(null), extraFilesInput = ref(null), replacementInput = ref(null)
const visible = ref(false), busy = ref(false), reading = ref(false), error = ref('')
const files = shallowRef([]), markdownFile = shallowRef(null)
const selectedPath = ref(''), source = ref(null), images = ref([])
const markdowns = computed(() => files.value.filter(file => /\.(?:md|markdown)$/i.test(file.name)))
const missingCount = computed(() => images.value.filter(image => !image.file).length)
const cache = new Map()
let replacementIndex = null, readVersion = 0
const filePath = file => file.webkitRelativePath || file.name
const isImage = file => /\.(?:png|jpe?g|gif|webp)$/i.test(file.name)

defineExpose({ selectMarkdown: () => markdownInput.value.click(), selectDirectory: () => directoryInput.value.click() })

function readText(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result))
    reader.onerror = () => reject(new Error('文章读取失败，请重新选择文件'))
    reader.readAsText(file)
  })
}

async function receiveFiles(event, directory) {
  const chosen = Array.from(event.target.files || [])
  event.target.value = ''
  if (!chosen.length) return
  files.value = chosen
  cache.clear()
  error.value = ''
  source.value = null
  images.value = []
  markdownFile.value = null
  selectedPath.value = ''
  visible.value = true
  readVersion++
  if (!markdowns.value.length) { error.value = '所选文件夹中没有 .md 或 .markdown 文件'; return }
  if (markdowns.value.length === 1) {
    selectedPath.value = filePath(markdowns.value[0])
    await loadSelected(!directory)
  }
}

async function loadSelected(autoEnter = false) {
  const version = ++readVersion
  const file = markdowns.value.find(file => filePath(file) === selectedPath.value)
  markdownFile.value = file || null
  source.value = null
  images.value = []
  error.value = ''
  if (!file) return
  reading.value = true
  try {
    const content = await readText(file)
    const { scanLocalImages, matchLocalImage } = await import('../utils/markdownImport')
    if (version !== readVersion) return
    source.value = content
    images.value = scanLocalImages(content).map(image => ({ ...image, ...matchLocalImage(image.source, file, files.value.filter(isImage)), status: 'ready', error: '' }))
    if (autoEnter === true && !images.value.length) await finish()
  } catch (err) {
    if (version === readVersion) { source.value = null; error.value = err.message }
  } finally {
    if (version === readVersion) reading.value = false
  }
}

async function addFiles(event) {
  const { matchLocalImage } = await import('../utils/markdownImport')
  files.value = [...files.value, ...Array.from(event.target.files || []).filter(isImage)]
  event.target.value = ''
  for (const image of images.value) {
    const match = matchLocalImage(image.source, markdownFile.value, files.value.filter(isImage))
    image.candidates = match.candidates
    if (!image.file && match.file) image.file = match.file
  }
}

function optionsFor(image) {
  return image.file && !image.candidates.includes(image.file) ? [image.file, ...image.candidates] : image.candidates
}
function chooseFile(image, event) {
  image.file = event.target.value === '' ? null : files.value[Number(event.target.value)]
  image.error = ''
  image.status = 'ready'
}
function selectReplacement(index) { replacementIndex = index; replacementInput.value.click() }
function replaceFile(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file || replacementIndex === null) return
  files.value = [...files.value, file]
  Object.assign(images.value[replacementIndex], { file, status: 'ready', error: '' })
}
function statusLabel(image) {
  if (image.status === 'uploaded') return '已上传'
  if (image.status === 'uploading') return '上传中…'
  return image.file ? '已匹配，待上传' : '未找到唯一匹配，请补选或选择图片'
}
async function finish() {
  if (busy.value || !markdownFile.value || source.value === null || missingCount.value) return
  busy.value = true
  error.value = ''
  try {
    const { importImages } = await import('../utils/markdownImport')
    const content = await importImages({ content: source.value, images: images.value, upload: uploadImageFile, cache })
    const title = content.match(/^#\s+(.+)\r?$/m)?.[1]?.trim() || markdownFile.value.name.replace(/\.(?:md|markdown)$/i, '')
    emit('imported', { content, title })
    visible.value = false
  } catch (err) {
    error.value = err.response?.data?.error || err.message || '导入失败，请重试'
  } finally { busy.value = false }
}
</script>

<style scoped>
.import-help { color: var(--muted); line-height: 1.7; margin: 0 0 16px; }
.import-filename { overflow-wrap: anywhere; margin-bottom: 16px; }
.article-choice { display: grid; gap: 8px; margin-bottom: 16px; }
select { width: 100%; min-width: 0; color: var(--fg); background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); padding: 8px; font: inherit; }
.import-toolbar { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 12px; margin: 16px 0; }
.image-list { max-height: 45vh; overflow: auto; margin-bottom: 16px; }
.image-row { padding: 14px 0; border-bottom: 1px solid var(--border); }
.image-source { display: grid; gap: 6px; margin-bottom: 10px; overflow-wrap: anywhere; }
.image-source span { color: var(--muted); font-size: 12px; }
.image-source .image-error { color: var(--error); }
.image-match { display: flex; align-items: center; gap: 8px; }
.image-match select { flex: 1; }
</style>
