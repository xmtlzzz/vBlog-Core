import { ref, onUnmounted } from 'vue'

export function useCountUp(getTarget, { duration = 2000 } = {}) {
  const display = ref(0)
  let animFrame = null
  let observer = null
  let inView = false
  let started = false

  function animate(final) {
    if (animFrame) cancelAnimationFrame(animFrame)
    if (final <= 0) { display.value = 0; return }
    const start = performance.now()

    function step(now) {
      const progress = Math.min((now - start) / duration, 1)
      const ease = 1 - Math.pow(1 - progress, 3)
      display.value = Math.round(final * ease)
      if (progress < 1) {
        animFrame = requestAnimationFrame(step)
      }
    }
    animFrame = requestAnimationFrame(step)
  }

  function mount(node) {
    observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        inView = true
        observer.disconnect()
        tryStart()
      }
    }, { threshold: 0.3 })
    observer.observe(node)
  }

  function tryStart() {
    if (started || !inView) return
    const val = getTarget()
    if (val > 0) {
      started = true
      animate(val)
    }
  }

  function start() {
    tryStart()
  }

  onUnmounted(() => {
    if (animFrame) cancelAnimationFrame(animFrame)
    if (observer) observer.disconnect()
  })

  return { display, mount, start }
}
