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
  effect: { type: String, default: '' }
})

const emit = defineEmits(['viewchange'])

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
      props.effect ? { effect: props.effect } : {},
      props.model === 'terrain' ? 'terrain' : 'tree',
      { onError: () => { if (rev === revision) showFallback(qr) } }
    )
    renderer.setFlat(flat.value)
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
  renderer?.setFlat(flat.value)
  emit('viewchange', flat.value)
}

watch([() => props.url, () => props.model], render)
watch(
  () => props.effect,
  (v) => renderer?.setScene(v ? { effect: v } : {})
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
