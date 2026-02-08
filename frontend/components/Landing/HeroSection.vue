<template>
  <div class="w-full max-w-7xl mx-auto px-4 py-12 lg:py-20">
    <div class="flex flex-col lg:grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
      
      <!-- LEFT SIDE: Brand Illustration -->
      <div class="order-2 lg:order-1 hidden lg:flex flex-col items-center justify-center">
        <!-- Interactive Clock Logo -->
        <div class="relative w-72 h-72 xl:w-96 xl:h-96">
          <!-- Pulsing Background Orbs (subtle, since main background is in index.vue) -->
          <div class="absolute inset-0">
            <div class="absolute inset-0 rounded-full bg-gradient-to-br from-primary-500/10 to-secondary-500/10 animate-pulse-slow blur-xl"></div>
            <div class="absolute inset-4 rounded-full bg-gradient-to-tr from-primary-500/5 to-secondary-500/5 animate-pulse-delayed blur-lg"></div>
          </div>

          <!-- Clock SVG -->
          <svg 
            class="w-full h-full drop-shadow-2xl"
            viewBox="0 0 200 200" 
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <!-- Clock Face Background -->
            <defs>
              <linearGradient id="clockFace" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#1e1b4b" />
                <stop offset="100%" stop-color="#312e81" />
              </linearGradient>
              <linearGradient id="clockRing" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#6366f1" />
                <stop offset="100%" stop-color="#a855f7" />
              </linearGradient>
              <filter id="glow">
                <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
                <feMerge>
                  <feMergeNode in="coloredBlur"/>
                  <feMergeNode in="SourceGraphic"/>
                </feMerge>
              </filter>
            </defs>

            <!-- Outer Ring -->
            <circle cx="100" cy="100" r="95" stroke="url(#clockRing)" stroke-width="4" fill="none" class="opacity-50" />
            
            <!-- Inner Clock Face -->
            <circle cx="100" cy="100" r="85" fill="url(#clockFace)" stroke="rgba(99, 102, 241, 0.3)" stroke-width="1" />

            <!-- Hour Markers -->
            <g class="clock-markers">
              <line x1="100" y1="25" x2="100" y2="35" stroke="#6366f1" stroke-width="3" stroke-linecap="round" />
              <line x1="100" y1="165" x2="100" y2="175" stroke="#6366f1" stroke-width="3" stroke-linecap="round" />
              <line x1="25" y1="100" x2="35" y2="100" stroke="#6366f1" stroke-width="3" stroke-linecap="round" />
              <line x1="165" y1="100" x2="175" y2="100" stroke="#6366f1" stroke-width="3" stroke-linecap="round" />
              
              <line x1="47" y1="47" x2="54" y2="54" stroke="#a855f7" stroke-width="2" stroke-linecap="round" opacity="0.7" />
              <line x1="146" y1="47" x2="153" y2="54" stroke="#a855f7" stroke-width="2" stroke-linecap="round" opacity="0.7" />
              <line x1="47" y1="153" x2="54" y2="146" stroke="#a855f7" stroke-width="2" stroke-linecap="round" opacity="0.7" />
              <line x1="146" y1="153" x2="153" y2="146" stroke="#a855f7" stroke-width="2" stroke-linecap="round" opacity="0.7" />
            </g>

            <!-- Center Dot -->
            <circle cx="100" cy="100" r="6" fill="url(#clockRing)" filter="url(#glow)" />
            <circle cx="100" cy="100" r="3" fill="#ffffff" />

            <!-- Hour Hand -->
            <g :style="{ transform: `rotate(${hourRotation}deg)`, transformOrigin: '100px 100px', transition: 'transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1)' }">
              <line x1="100" y1="100" x2="100" y2="55" stroke="#ffffff" stroke-width="5" stroke-linecap="round" filter="url(#glow)" />
              <line x1="100" y1="100" x2="100" y2="55" stroke="#e0e7ff" stroke-width="3" stroke-linecap="round" />
            </g>

            <!-- Minute Hand -->
            <g :style="{ transform: `rotate(${minuteRotation}deg)`, transformOrigin: '100px 100px', transition: 'transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)' }">
              <line x1="100" y1="100" x2="100" y2="35" stroke="#a855f7" stroke-width="4" stroke-linecap="round" filter="url(#glow)" />
              <line x1="100" y1="100" x2="100" y2="35" stroke="#c4b5fd" stroke-width="2" stroke-linecap="round" />
            </g>

            <!-- Second Hand -->
            <g :style="{ transform: `rotate(${secondRotation}deg)`, transformOrigin: '100px 100px', transition: 'transform 0.1s ease-out' }">
              <line x1="100" y1="115" x2="100" y2="30" stroke="#06c755" stroke-width="2" stroke-linecap="round" opacity="0.9" />
            </g>

            <!-- Brand Text -->
            <text x="100" y="195" text-anchor="middle" class="font-sans text-xs fill-slate-400 tracking-widest">TIMELEDGER</text>
          </svg>
        </div>

        <!-- Brand Tagline -->
        <div class="mt-8 text-center">
          <h2 class="font-sans text-2xl xl:text-3xl font-semibold bg-gradient-to-r from-white via-indigo-200 to-white bg-clip-text text-transparent">
            教師排課首選平台
          </h2>
          <p class="mt-2 text-slate-400 font-sans text-sm xl:text-base">
            智慧排課 · 人才媒合 · 雲端管理
          </p>
        </div>
      </div>

      <!-- RIGHT SIDE: Login Card -->
      <div class="order-1 lg:order-2 w-full max-w-md mx-auto lg:mx-0">
        <!-- Mobile Logo (visible only on mobile) -->
        <div class="lg:hidden mb-6 text-center">
          <div class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-primary-500 to-secondary-500 shadow-xl shadow-primary-500/20">
            <svg class="w-10 h-10 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3-4 4m4-4h6a2 2 0 012 2v8a2 2 0 01-2 2H6a2 2 0 00-2 2v12" />
            </svg>
          </div>
        </div>

        <!-- Card Header -->
        <div class="text-center mb-8">
          <h1 class="font-heading text-3xl font-bold text-white mb-2 animate-slide-up-1">
            歡迎回來
          </h1>
          <p class="font-sans text-slate-400 text-sm animate-slide-up-2">
            登入您的教師帳號
          </p>
        </div>

        <!-- LINE Login Button (Teacher Focus) -->
        <button
          @click="$emit('line-login')"
          :disabled="loading"
          class="w-full line-login-btn flex items-center justify-center gap-3 py-4 px-6 rounded-xl font-sans text-base font-medium transition-all duration-300 animate-slide-up-3"
        >
          <svg v-if="loading" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v16a8 8 0 00-8 8h16a8 8 0 018-8z" />
          </svg>
          <svg v-else class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.477 2 2 6.477 2 12c0 1.891.5 3.657 1.363 5.203L2 22l4.797-1.363A9.956 9.956 0 0 0 12 22c5.523 0 10-4.477 10-10S17.523 2 12 2z"/>
          </svg>
          <span>以 LINE 快速登入</span>
        </button>

        <!-- Divider -->
        <div class="relative my-6 animate-slide-up-4">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-white/10"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-4 bg-transparent text-slate-500 font-sans text-xs uppercase tracking-wider">
              或使用管理員帳號
            </span>
          </div>
        </div>

        <!-- Admin Login Section -->
        <div class="animate-slide-up-4">
          <NuxtLink
            to="/admin/login"
            class="admin-entry-btn group flex items-center justify-between w-full py-4 px-5 rounded-xl glass-card transition-all duration-300"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-indigo-500/20 to-purple-500/20 flex items-center justify-center group-hover:from-indigo-500/30 group-hover:to-purple-500/30 transition-all duration-300">
                <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
              </div>
              <div class="text-left">
                <p class="font-heading text-white font-medium text-sm">管理員登入</p>
                <p class="font-sans text-slate-500 text-xs">中心管理人員專用</p>
              </div>
            </div>
            <svg class="w-5 h-5 text-slate-500 group-hover:text-primary-400 transform group-hover:translate-x-1 transition-all duration-300" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
            </svg>
          </NuxtLink>

          <!-- Admin Quick Link -->
          <p class="mt-4 text-center">
            <NuxtLink to="/admin/login" class="font-sans text-xs text-slate-500 hover:text-primary-400 transition-colors duration-200 underline-offset-2 hover:underline">
              前往管理員後台
            </NuxtLink>
          </p>
        </div>

        <!-- Feature Highlights -->
        <div class="mt-8 pt-6 border-t border-white/5 animate-fade-in-5">
          <div class="grid grid-cols-3 gap-4 text-center">
            <div class="feature-item">
              <div class="text-primary-400 font-heading font-semibold text-lg">98%</div>
              <div class="text-slate-500 font-sans text-xs mt-1">排課效率</div>
            </div>
            <div class="feature-item">
              <div class="text-primary-400 font-heading font-semibold text-lg">500+</div>
              <div class="text-slate-500 font-sans text-xs mt-1">合作中心</div>
            </div>
            <div class="feature-item">
              <div class="text-primary-400 font-heading font-semibold text-lg">10K+</div>
              <div class="text-slate-500 font-sans text-xs mt-1">教師信賴</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

defineEmits(['line-login'])

defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
})

// Clock animation state
const currentTime = ref(new Date())
let timeInterval: ReturnType<typeof setInterval> | null = null

// Calculate clock hand rotations
const secondRotation = computed(() => {
  return currentTime.value.getSeconds() * 6
})

const minuteRotation = computed(() => {
  return currentTime.value.getMinutes() * 6 + currentTime.value.getSeconds() * 0.1
})

const hourRotation = computed(() => {
  const hours = currentTime.value.getHours() % 12
  const minutes = currentTime.value.getMinutes()
  return hours * 30 + minutes * 0.5
})

// Update time every second
const updateTime = () => {
  currentTime.value = new Date()
}

onMounted(() => {
  timeInterval = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
})
</script>

<style scoped>
/* Glass Card - Unified styling */
.glass-card {
  @apply bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl shadow-black/20;
}

.glass-card:hover {
  @apply border-white/20;
  box-shadow: 0 25px 50px -12px rgba(99, 102, 241, 0.15);
}

/* Pulse animations */
.animate-pulse-slow {
  animation: pulse-slow 4s ease-in-out infinite;
}

.animate-pulse-delayed {
  animation: pulse-slow 4s ease-in-out infinite 2s;
}

@keyframes pulse-slow {
  0%, 100% {
    opacity: 0.5;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}

/* LINE Login Button */
.line-login-btn {
  @apply bg-[#06C755] text-white shadow-lg shadow-green-500/25;
  @apply hover:bg-[#05b547] hover:shadow-green-500/40;
  @apply active:scale-[0.98] hover:-translate-y-0.5;
  @apply hover:scale-[1.02] transition-all duration-300;
}

/* Admin Entry Button */
.admin-entry-btn {
  @apply hover:bg-slate-800/60;
}

.admin-entry-btn:hover {
  @apply border-primary-500/30;
}

/* Page Load Animations */
@keyframes slide-up-1 {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slide-up-2 {
  0% {
    opacity: 0;
    transform: translateY(20px);
    animation-delay: 0.1s;
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slide-up-3 {
  0% {
    opacity: 0;
    transform: translateY(20px);
    animation-delay: 0.2s;
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slide-up-4 {
  0% {
    opacity: 0;
    transform: translateY(20px);
    animation-delay: 0.3s;
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-in-5 {
  from {
    opacity: 0;
    animation-delay: 0.5s;
  }
  to {
    opacity: 1;
  }
}

.animate-slide-up-1 {
  animation: slide-up-1 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.animate-slide-up-2 {
  animation: slide-up-2 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.animate-slide-up-3 {
  animation: slide-up-3 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.animate-slide-up-4 {
  animation: slide-up-4 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.animate-fade-in-5 {
  animation: fade-in-5 0.5s ease-out forwards;
}

/* Feature items */
.feature-item {
  @apply transition-all duration-300 hover:scale-105;
}

/* Responsive adjustments */
@media (max-width: 1024px) {
  .glass-card {
    @apply p-6;
  }
}

@media (max-width: 640px) {
  .glass-card {
    @apply p-5 rounded-2xl;
  }

  .line-login-btn {
    @apply py-3.5;
  }
}
</style>
