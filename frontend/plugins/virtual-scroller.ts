import VueVirtualScroller from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

// 追蹤是否已安裝
let installed = false

export default defineNuxtPlugin((nuxtApp) => {
  if (!installed) {
    nuxtApp.vueApp.use(VueVirtualScroller)
    installed = true
  }
})
