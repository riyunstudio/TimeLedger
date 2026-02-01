<template>
  <header class="relative overflow-hidden">
    <div class="absolute inset-0 z-0">
      <!-- Animated Mesh Gradient Background -->
      <div class="absolute inset-0 animate-gradient">
        <div class="gradient-orb-1"></div>
        <div class="gradient-orb-2"></div>
        <div class="gradient-orb-3"></div>
        <div class="gradient-orb-4"></div>
        <div class="gradient-orb-5"></div>
        <div class="gradient-orb-6"></div>
      </div>
    </div>

    <!-- Content Overlay -->
    <div class="relative z-10">
      <div class="max-w-6xl mx-auto px-4 py-20 text-center">
        <!-- Logo -->
        <div class="mb-8">
          <div class="w-24 h-24 mx-auto mb-4 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-full flex items-center justify-center shadow-2xl">
            <svg class="w-12 h-12 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-6-4 4m6-6h6a2 2 0 012 0V8a2 2 0 01-2-2H6a2 2 0 00-2 2v12" />
            </svg>
          </div>
        </div>

        <!-- Title -->
        <h1 class="text-5xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-show-in-view-fade mb-4">
          {{ $t('landing.hero.title') }}
        </h1>
        <p class="text-2xl text-slate-300 font-medium">
          {{ $t('landing.hero.subtitle') }}
        </p>

        <!-- Slogan -->
        <div class="mb-12">
          <p class="text-xl text-white font-medium">
            {{ $t('landing.hero.slogan_multi_center') }}
            <span class="mx-4 text-slate-400">× </span>
            {{ $t('landing.hero.slogan_talent_pool') }}
            <span class="text-white/80">× </span>
            {{ $t('landing.hero.slogan_schedule_share') }}
          </p>
        </div>

        <!-- CTA Buttons -->
        <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
          <button
            @click="$emit('line-login')"
            :disabled="loading"
            class="w-full sm:w-auto btn-primary flex items-center justify-center gap-3 py-4 px-8 text-lg shadow-2xl"
          >
            <svg v-if="loading" class="animate-spin w-6 h-6" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v16a8 8 0 00-8 8h16a8 8 0 018-8zm-9 8a8 8 0 11-14 0v16a8 8 0 00-16 8H4a8 8 0 01-8 8v-16zm10 0a6 6 0 00-6 6v20a6 6 0 016 6H4a6 6 0 0016 6v20zM2 2a8 8 0 01-8 8v24a8 8 0 018-8 8H4a8 8 0 0016 8V6z" />
            </svg>
            <svg v-else class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2C6.477 2 2 6.477 2 12c0 1.891.5 3.657 1.363 5.203L2 22l4.797-1.363A9.956 9.956 0 0 0 12 22c5.523 0 10-4.477 10-10S17.523 2 12 2zm0 18c-1.455 0-2.822-.395-4-1.088V16c0-1.103.897-2 2-2h4c1.103 0 2 .897 2 2v.912C14.822 19.605 13.455 20 12 20zm6-5h-1v-1c0-2.206-1.794-4-4-4s-4 1.794-4 4v1H8v-1c0-2.762 2.238-5 5-5s5 2.238 5 5v1z"/>
            </svg>
            <span class="hidden sm:inline">{{ $t('landing.hero.login_line_full') }}</span>
            <span class="sm:hidden">{{ $t('landing.hero.login_line') }}</span>
          </button>
          <NuxtLink
            to="/admin/login"
            class="w-full sm:w-auto glass-btn px-8 py-4 text-lg font-medium hover:bg-white/10 transition-all"
          >
            {{ $t('landing.hero.login_admin') }}
          </NuxtLink>
        </div>

        <!-- Admin Entry Link -->
        <p class="mt-4">
          <NuxtLink to="/admin/login" class="text-sm text-slate-400 underline-offset-2 hover:text-primary-500 transition-colors">
            {{ $t('landing.hero.login_admin_hint') }}
          </NuxtLink>
        </p>
      </div>
    </div>
  </header>
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
.animate-gradient {
  animation: gradient-shift 20s linear infinite;
}

@keyframes gradient-shift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.gradient-orb-1,
.gradient-orb-2,
.gradient-orb-3,
.gradient-orb-4,
.gradient-orb-5,
.gradient-orb-6 {
  position: absolute;
  border-radius: 50%;
  animation: gradient-shift 20s linear infinite;
}

.gradient-orb-1 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #6366F1 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-2 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #A855F7 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-3 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #1e3a8a 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-4 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #1e3a8a 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-5 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #1e3a8a 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-6 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #6366F1 0%, transparent 50%);
  background-size: 200% 200%;
}

.bg-show-in-view-fade {
  animation: show-in-fade 1s ease-out forwards;
}

@keyframes show-in-fade {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
