<template>
  <button
    type="button"
    class="eqr"
    :aria-label="fallback ? 'QR code' : flat ? '恢复模型视图' : '显示二维码'"
    @click="toggleView"
  >
    <canvas v-show="!fallback" ref="canvasRef" :key="canvasKey"></canvas>
    <svg
      v-if="fallback"
      class="eqr-fallback"
      :viewBox="`0 0 ${fallback.size} ${fallback.size}`"
      shape-rendering="crispEdges"
      aria-hidden="true"
    >
      <rect :width="fallback.size" :height="fallback.size" fill="#ffffff" />
      <path :d="fallback.path" fill="#111111" />
    </svg>
  </button>
</template>

<script setup>
import { ref, watch, nextTick, onBeforeUnmount } from 'vue'
import {
  CURRENT_GENERATOR_VERSION,
  createEveryQRCodeIdentity,
  createQRSvgPath,
} from '@every-qrcode/core'

const props = defineProps({
  url: { type: String, required: true },
  // tree（樱花树）| terrain（地形）
  model: { type: String, default: 'tree' },
  // ''（库默认）| calm | snow | rain | wind
  effect: { type: String, default: '' },
  // sakura（库默认樱花）| forest | ocean | sunset（仅作用于树模型）
  palette: { type: String, default: 'sakura' }
})

const emit = defineEmits(['viewchange'])

// 五色槽含义（对齐库默认樱花配色）：[花, 瓣, 亮块, 枝干/二维码模块, 草地]
const QR_PALETTES = {
  sakura: null,
  forest: [
    [0.55, 0.76, 0.45],
    [0.72, 0.85, 0.56],
    [0.88, 0.91, 0.79],
    [0.25, 0.38, 0.16],
    [0.92, 0.94, 0.86]
  ],
  ocean: [
    [0.45, 0.72, 0.86],
    [0.63, 0.83, 0.91],
    [0.86, 0.93, 0.96],
    [0.15, 0.32, 0.46],
    [0.88, 0.94, 0.96]
  ],
  sunset: [
    [0.96, 0.55, 0.32],
    [1.0, 0.76, 0.42],
    [0.96, 0.87, 0.72],
    [0.5, 0.22, 0.12],
    [0.99, 0.92, 0.82]
  ]
}

function sceneConfig() {
  const scene = {}
  if (props.effect) scene.effect = props.effect
  const palette = QR_PALETTES[props.palette]
  if (palette) scene.palette = palette
  return scene
}

// 模型视图放大让树更饱满（库钳制 0.82–1.45）；超过约 1.25 会裁掉树冠，
// 二维码视图保持 1.0，避免裁掉静区
const TREE_ZOOM = 1.2

function applyView(renderer) {
  renderer.setFlat(flat.value)
  renderer.setZoom(flat.value ? 1 : TREE_ZOOM)
}

const canvasRef = ref(null)
const canvasKey = ref(0)
const fallback = ref(null)
const flat = ref(false)

let renderer = null
let observer = null
let revision = 0

function showFallback(qr) {
  renderer?.dispose()
  renderer = null
  fallback.value = qr
  emit('viewchange', true)
}

function disposeRenderer() {
  observer?.disconnect()
  observer = null
  renderer?.dispose()
  renderer = null
}

async function render() {
  const rev = ++revision
  disposeRenderer()
  const identity = await createEveryQRCodeIdentity(props.url, { identityScope: 'url' }).catch(() => null)
  if (rev !== revision) return
  if (!identity) return
  const qr = createQRSvgPath(identity.qr)
  fallback.value = null
  // WebGPU 上下文不能在同一个 canvas 上重复获取，重渲染时换新 canvas
  canvasKey.value++
  await nextTick()
  const canvas = canvasRef.value
  if (!canvas) return
  try {
    const { createSeedModel, mountSeed } = await import('@every-qrcode/renderer-webgpu')
    const seed = await createSeedModel(identity, { generatorVersion: CURRENT_GENERATOR_VERSION })
    if (rev !== revision) return
    renderer = mountSeed(
      canvas,
      seed,
      sceneConfig(),
      props.model === 'terrain' ? 'terrain' : 'tree',
      { onError: () => { if (rev === revision) showFallback(qr) } }
    )
    applyView(renderer)
    renderer.resize()
    observer = new ResizeObserver(() => renderer?.resize())
    observer.observe(canvas)
  } catch {
    if (rev === revision) showFallback(qr)
  }
}

function toggleView() {
  if (fallback.value) return
  flat.value = !flat.value
  if (renderer) applyView(renderer)
  emit('viewchange', flat.value)
}

watch([() => props.url, () => props.model], render)
watch(
  [() => props.effect, () => props.palette],
  () => renderer?.setScene(sceneConfig())
)

render()

onBeforeUnmount(() => {
  revision++
  disposeRenderer()
})
</script>

<style scoped>
.eqr {
  display: block;
  width: 100%;
  height: 100%;
  padding: 0;
  border: 0;
  background: transparent;
  cursor: pointer;
}
.eqr canvas,
.eqr svg {
  display: block;
  width: 100%;
  height: 100%;
}
</style>
