<template>
  <div class="max-w-5xl mx-auto mt-20 px-4">
    <h2 class="text-3xl font-bold text-slate-100 text-center mb-8">
      {{ $t('landing.features.title') }}
    </h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Feature 1: Beautiful Export -->
      <div class="glass-card p-6 rounded-2xl">
        <div class="aspect-[9/16] rounded-xl overflow-hidden mb-4 relative">
          <div class="absolute inset-0 bg-gradient-mesh flex items-center justify-center">
            <div class="text-center">
              <h3 class="text-2xl font-bold text-white mb-2">{{ $t('landing.features.export.sample_title') }}</h3>
              <div class="flex items-center justify-center gap-2 mb-4">
                <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
                  <span class="text-xl font-bold text-white">A</span>
                </div>
                <div class="space-y-2">
                  <div class="bg-white/20 rounded-lg p-2">
                    <p class="text-xs text-slate-800">14:00 鋼琴基礎</p>
                    <p class="text-xs text-slate-800">15:00 小提琴入門</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <h3 class="text-lg font-semibold text-slate-100 text-center">{{ $t('landing.features.export.title') }}</h3>
        <p class="text-sm text-slate-400 text-center mt-2">
          {{ $t('landing.features.export.description') }}
        </p>
      </div>

      <!-- Feature 2: Smart Matching -->
      <div class="glass-card p-6 rounded-2xl">
        <div class="p-4 mb-4 rounded-xl bg-secondary-500/20 border border-secondary-500/50 flex items-center gap-4">
          <svg class="w-10 h-10 text-secondary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4" />
          </svg>
          <div>
            <h4 class="font-semibold text-slate-100">{{ $t('landing.features.matching.title') }}</h4>
            <p class="text-sm text-slate-400 mt-2">
              {{ $t('landing.features.matching.description') }}
            </p>
          </div>
        </div>
        <div class="mt-4">
          <div class="bg-slate-800/50 rounded-lg p-3">
            <p class="text-xs text-white mb-1">🎯 Alice</p>
            <p class="text-xs text-slate-200">鋼琴 · 高級 · 5星</p>
          </div>
        </div>
        <div class="flex items-center gap-4">
          <div class="bg-primary-500/20 rounded-lg px-3 py-2">
            <p class="text-xs font-white font-medium">技能</p>
            <span class="px-2 py-0.5 rounded-full bg-white/20 text-xs text-primary-500">鋼琴</span>
            <span class="px-2 py-0.5 rounded-full bg-white/20 text-xs text-secondary-500">古典</span>
          </div>
          <div class="bg-warning-500/20 rounded-lg px-3 py-2">
            <p class="text-xs text-white font-medium">證照</p>
            <span class="px-2 py-0.5 rounded-full bg-success-500/20 text-xs text-white">ABRSM</span>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <p class="text-xs text-slate-200">96 分匹配</p>
        </div>
      </div>

      <!-- Feature 3: Real-time Notifications -->
      <div class="glass-card p-6 rounded-2xl">
        <div class="p-4 mb-4 rounded-xl bg-critical-500/20 border border-critical-500/50 flex items-center gap-4">
          <svg class="w-10 h-10 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5v-2a3 3 0 00-3-3H6a3 3 0 00-3 3v11a3 3 0 003 3h2a2 2 0 002-2v-5a2 2 0 00-2-2h-1" />
          </svg>
          <div>
            <h4 class="font-semibold text-slate-100">{{ $t('landing.features.notifications.title') }}</h4>
            <p class="text-sm text-slate-400 mt-2">
              {{ $t('landing.features.notifications.description') }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- CTA -->
    <div class="text-center mt-12">
      <button
        @click="$emit('line-login')"
        :disabled="loading"
        class="btn-primary px-8 py-4 text-lg shadow-2xl animate-bounce"
      >
        {{ loading ? $t('landing.cta.connecting') : $t('landing.cta.start_using') }}
      </button>
    </div>

    <!-- Footer -->
    <div class="text-center mt-20 text-slate-500 text-sm">
      <p>{{ $t('landing.footer.copyright') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
defineEmits(['line-login'])

defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
})
</script>

<style scoped>
.animate-bounce {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-20px);
  }
  60% {
    transform: translateY(-10px);
  }
}
</style>
