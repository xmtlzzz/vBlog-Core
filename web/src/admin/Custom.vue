<template>
  <div class="custom-page">
    <div class="page-header">
      <h1 class="page-title">组件定制</h1>
      <el-button type="primary" @click="openUpload">上传组件 Upload</el-button>
    </div>

    <!-- Layout components -->
    <div class="section" v-if="layoutComponents.length">
      <h2 class="section-title">Layout</h2>
      <div class="comp-grid fade-in stagger">
        <div v-for="comp in layoutComponents" :key="comp.id" class="comp-card card-hover">
          <div class="comp-header">
            <span class="comp-name">{{ comp.name }}</span>
            <span :class="['status-badge', comp.status === 'active' ? 'active' : 'inactive']">
              {{ comp.status === 'active' ? '启用' : '停用' }}
            </span>
          </div>
          <p class="comp-desc">{{ comp.description || '无描述' }}</p>
          <div class="comp-meta">
            <span>v{{ comp.version || '1.0' }}</span>
            <span>{{ comp.origin || 'built-in' }}</span>
          </div>
          <div class="comp-actions">
            <el-button size="small" :type="comp.status === 'active' ? 'warning' : 'success'" text @click="toggleComponent(comp)">
              {{ comp.status === 'active' ? '停用' : '启用' }}
            </el-button>
            <el-button v-if="comp.origin !== 'built-in'" size="small" type="danger" text @click="deleteComponent(comp.id)">删除</el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- Content components -->
    <div class="section" v-if="contentComponents.length">
      <h2 class="section-title">Content</h2>
      <div class="comp-grid fade-in stagger">
        <div v-for="comp in contentComponents" :key="comp.id" class="comp-card card-hover">
          <div class="comp-header">
            <span class="comp-name">{{ comp.name }}</span>
            <span :class="['status-badge', comp.status === 'active' ? 'active' : 'inactive']">
              {{ comp.status === 'active' ? '启用' : '停用' }}
            </span>
          </div>
          <p class="comp-desc">{{ comp.description || '无描述' }}</p>
          <div class="comp-meta">
            <span>v{{ comp.version || '1.0' }}</span>
            <span>{{ comp.origin || 'built-in' }}</span>
          </div>
          <div class="comp-actions">
            <el-button size="small" :type="comp.status === 'active' ? 'warning' : 'success'" text @click="toggleComponent(comp)">
              {{ comp.status === 'active' ? '停用' : '启用' }}
            </el-button>
            <el-button v-if="comp.origin !== 'built-in'" size="small" type="danger" text @click="deleteComponent(comp.id)">删除</el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- Custom uploads -->
    <div class="section" v-if="customComponents.length">
      <h2 class="section-title">Custom Uploads</h2>
      <div class="comp-grid fade-in stagger">
        <div v-for="comp in customComponents" :key="comp.id" class="comp-card card-hover">
          <div class="comp-header">
            <span class="comp-name">{{ comp.name }}</span>
            <span :class="['status-badge', comp.status === 'active' ? 'active' : 'inactive']">
              {{ comp.status === 'active' ? '启用' : '停用' }}
            </span>
          </div>
          <p class="comp-desc">{{ comp.description || '无描述' }}</p>
          <div class="comp-meta">
            <span>v{{ comp.version || '1.0' }}</span>
            <span>{{ comp.origin || 'custom' }}</span>
          </div>
          <div class="comp-actions">
            <el-button size="small" type="primary" text @click="previewComponent(comp)">预览</el-button>
            <el-button size="small" :type="comp.status === 'active' ? 'warning' : 'success'" text @click="toggleComponent(comp)">
              {{ comp.status === 'active' ? '停用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" text @click="deleteComponent(comp.id)">删除</el-button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="!components.length" class="empty-state">暂无组件</div>

    <!-- Preview Modal -->
    <el-dialog v-model="previewVisible" :title="`预览：${previewComp?.name || ''}`" width="680px" @closed="previewComp = null">
      <div class="preview-container">
        <iframe v-if="previewComp" :srcdoc="buildSrcdoc(previewComp.code)" class="preview-iframe" sandbox="allow-scripts"></iframe>
      </div>
    </el-dialog>

    <!-- Upload Modal -->
    <el-dialog v-model="dialogVisible" title="上传组件" width="560px" destroy-on-close>
      <el-form label-position="top">
        <el-form-item label="名称">
          <el-input v-model="form.name" placeholder="组件名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="组件描述" />
        </el-form-item>
        <el-form-item label="版本">
          <el-input v-model="form.version" placeholder="1.0.0" />
        </el-form-item>
        <el-form-item label="代码">
          <div class="code-toolbar">
            <el-upload :show-file-list="false" accept=".js,.jsx,.ts,.tsx" :before-upload="onImportFile">
              <el-button size="small" text type="primary">导入文件</el-button>
            </el-upload>
          </div>
          <el-input v-model="form.code" type="textarea" :rows="8" placeholder="组件代码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="uploadComponent">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'
import { buildSrcdoc as _buildSrcdoc } from '../utils/component'

const components = ref([])
const dialogVisible = ref(false)
const saving = ref(false)
const form = reactive({ name: '', description: '', version: '', code: '' })
const previewVisible = ref(false)
const previewComp = ref(null)

function buildSrcdoc(code) {
  if (!code) return '<body style="padding:24px;color:#737373;font-family:system-ui">该组件没有代码内容</body>'
  return _buildSrcdoc(code, { padding: '24px' })
}

const layoutComponents = computed(() => components.value.filter(c => c.category === 'layout'))
const contentComponents = computed(() => components.value.filter(c => c.category === 'content'))
const customComponents = computed(() => components.value.filter(c => !c.category || c.category === 'custom'))

async function fetchComponents() {
  const res = await api.get('/components').catch(() => ({ data: [] }))
  components.value = Array.isArray(res) ? res : (res.data || [])
}

function openUpload() {
  Object.assign(form, { name: '', description: '', version: '', code: '' })
  dialogVisible.value = true
}

function onImportFile(file) {
  const reader = new FileReader()
  reader.onload = (e) => {
    form.code = e.target.result
    if (!form.name) {
      form.name = file.name.replace(/\.(js|jsx|ts|tsx)$/, '')
    }
    ElMessage.success(`已导入 ${file.name}`)
  }
  reader.readAsText(file)
  return false
}

async function uploadComponent() {
  if (!form.name.trim()) {
    ElMessage.warning('请输入组件名称')
    return
  }
  saving.value = true
  try {
    await api.post('/components', {
      name: form.name,
      description: form.description,
      version: form.version,
      code: form.code,
      category: 'custom',
      origin: 'custom'
    })
    ElMessage.success('组件已上传')
    dialogVisible.value = false
    fetchComponents()
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

async function toggleComponent(comp) {
  await api.patch(`/components/${comp.id}/toggle`)
  comp.status = comp.status === 'active' ? 'inactive' : 'active'
  ElMessage.success(comp.status === 'active' ? '已启用' : '已停用')
}

async function deleteComponent(id) {
  try {
    await ElMessageBox.confirm('确定删除此组件？', '确认删除', { type: 'warning' })
    await api.delete(`/components/${id}`)
    ElMessage.success('已删除')
    fetchComponents()
  } catch {}
}

function previewComponent(comp) {
  previewComp.value = comp
  previewVisible.value = true
}

onMounted(fetchComponents)
</script>

<style scoped>
.custom-page {
  width: 100%;
}
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.page-title {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 600;
  color: var(--fg);
}
.section {
  margin-bottom: 28px;
}
.section-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--fg);
  margin-bottom: 12px;
}
.comp-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
.comp-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 20px;
  transition: border-color 0.15s;
}
.comp-card:hover {
  border-color: var(--accent);
}
.comp-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.comp-name {
  font-weight: 600;
  font-size: 15px;
  color: var(--fg);
}
.comp-desc {
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 12px;
  line-height: 1.5;
}
.comp-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--muted);
  margin-bottom: 12px;
}
.comp-actions {
  display: flex;
  gap: 8px;
}
.status-badge {
  display: inline-block;
  padding: 1px 8px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
}
.status-badge.active {
  background: var(--success-soft);
  color: var(--success);
}
.status-badge.inactive {
  background: var(--card-hover);
  color: var(--muted);
}
.empty-state {
  text-align: center;
  padding: 48px;
  color: var(--muted);
}
.code-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 4px;
}
.preview-container {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  min-height: 120px;
  padding: 16px;
}
.preview-iframe {
  width: 100%;
  min-height: 100px;
  border: none;
  background: transparent;
}
</style>
