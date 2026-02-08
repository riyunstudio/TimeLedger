<template>
  <div class="py-20 px-4">
    <!-- Section Header -->
    <div class="max-w-5xl mx-auto text-center mb-12">
      <h2 class="relative z-10 text-3xl lg:text-4xl font-bold text-slate-100 mb-4">
        {{ $t('landing.features.title') }}
      </h2>
      <p class="relative z-10 text-lg text-slate-400 max-w-2xl mx-auto">
        {{ $t('landing.features.subtitle') }}
      </p>
    </div>

    <div class="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
      <!-- Feature 1: Smart Scheduling -->
      <div
        class="group glass-card p-6 lg:p-8 rounded-2xl transition-all duration-500 hover:shadow-indigo-500/10 hover:-translate-y-2 opacity-0 animate-fade-in-up"
        style="animation-delay: 0ms;"
      >
        <!-- Icon Container -->
        <div class="relative mb-6">
          <div class="aspect-[9/16] rounded-xl overflow-hidden relative">
            <!-- Gradient background -->
            <div class="absolute inset-0 bg-gradient-to-br from-indigo-950/80 via-slate-900/90 to-slate-950"></div>

            <!-- Radial glow -->
            <div class="absolute inset-0 opacity-30" style="background-image: radial-gradient(circle at 50% 50%, rgba(99, 102, 241, 0.4) 0%, transparent 50%);"></div>

            <!-- Mini 7-Day Grid SVG with Drag-and-Drop Motion -->
            <div class="absolute inset-0 flex items-center justify-center p-3">
              <svg viewBox="0 0 80 80" class="w-full h-full drop-shadow-2xl" fill="none" xmlns="http://www.w3.org/2000/svg">
                <defs>
                  <linearGradient id="schedBlockGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#6366f1" />
                    <stop offset="100%" stop-color="#818cf8" />
                  </linearGradient>
                  <linearGradient id="schedBlockGrad2" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#4f46e5" />
                    <stop offset="100%" stop-color="#6366f1" />
                  </linearGradient>
                </defs>

                <!-- Calendar Grid Base -->
                <rect x="10" y="8" width="60" height="64" rx="4" fill="rgba(255,255,255,0.05)" stroke="rgba(99,102,241,0.3)" stroke-width="1" />

                <!-- Week Header -->
                <text x="20" y="20" font-size="5" fill="#6366f1" font-weight="600">一</text>
                <text x="30" y="20" font-size="5" fill="#6366f1" font-weight="600">二</text>
                <text x="40" y="20" font-size="5" fill="#6366f1" font-weight="600">三</text>
                <text x="50" y="20" font-size="5" fill="#6366f1" font-weight="600">四</text>
                <text x="60" y="20" font-size="5" fill="#6366f1" font-weight="600">五</text>

                <!-- Schedule Blocks with Drag Motion -->
                <!-- Block 1: Monday 10-12 -->
                <g class="animate-pulse" style="animation-duration: 3s;">
                  <rect x="18" y="26" width="10" height="8" rx="2" fill="url(#schedBlockGrad)" opacity="0.9" />
                  <rect x="19" y="27" width="8" height="2" rx="1" fill="rgba(255,255,255,0.3)" />
                </g>

                <!-- Block 2: Tuesday 14-16 (being dragged) -->
                <g>
                  <rect x="38" y="44" width="10" height="8" rx="2" fill="url(#schedBlockGrad2)" class="transition-all duration-500 group-hover:translate-x-1" opacity="0.9" />
                  <rect x="39" y="45" width="8" height="2" rx="1" fill="rgba(255,255,255,0.3)" />
                  <!-- Drag indicator -->
                  <circle cx="54" cy="48" r="3" fill="none" stroke="#6366f1" stroke-width="1" stroke-dasharray="2 1" class="animate-spin" style="animation-duration: 4s;" />
                </g>

                <!-- Block 3: Thursday 10-12 -->
                <g class="animate-pulse" style="animation-delay: 1s; animation-duration: 3s;">
                  <rect x="58" y="26" width="10" height="8" rx="2" fill="url(#schedBlockGrad)" opacity="0.8" />
                  <rect x="59" y="27" width="8" height="2" rx="1" fill="rgba(255,255,255,0.3)" />
                </g>

                <!-- Block 4: Friday 16-18 -->
                <g class="animate-pulse" style="animation-delay: 2s; animation-duration: 3s;">
                  <rect x="68" y="52" width="10" height="8" rx="2" fill="url(#schedBlockGrad)" opacity="0.7" />
                  <rect x="69" y="53" width="8" height="2" rx="1" fill="rgba(255,255,255,0.3)" />
                </g>

                <!-- Saturday indicator -->
                <text x="28" y="70" font-size="5" fill="#a855f7" font-weight="500">週六</text>
                <text x="48" y="70" font-size="5" fill="#a855f7" font-weight="500">週日</text>

                <!-- Conflict Detection Indicator -->
                <circle cx="14" cy="14" r="5" fill="#ef4444" class="animate-pulse">
                  <animate attributeName="opacity" values="1;0.5;1" dur="1.5s" repeatCount="indefinite" />
                </circle>
                <path d="M12 12 L14 14 M14 12 L12 14" stroke="white" stroke-width="1.5" stroke-linecap="round" />
              </svg>
            </div>

            <!-- Floating label -->
            <div class="absolute bottom-2 left-2 right-2 flex justify-center">
              <span class="text-[10px] text-indigo-300 bg-indigo-500/20 px-2 py-0.5 rounded-full border border-indigo-500/30">
                衝突偵測
              </span>
            </div>
          </div>
        </div>

        <!-- Content -->
        <h3 class="text-xl font-semibold text-slate-100 text-center mb-2">
          {{ $t('landing.features.scheduling.title') }}
        </h3>
        <p class="text-sm text-slate-400 text-center leading-relaxed">
          {{ $t('landing.features.scheduling.description') }}
        </p>

        <!-- Hover glow effect -->
        <div class="absolute inset-0 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"
             style="box-shadow: inset 0 0 40px rgba(99, 102, 241, 0.1);">
        </div>
      </div>

      <!-- Feature 2: Smart Matching -->
      <div
        class="group glass-card p-6 lg:p-8 rounded-2xl transition-all duration-500 hover:shadow-purple-500/10 hover:-translate-y-2 opacity-0 animate-fade-in-up"
        style="animation-delay: 150ms;"
      >
        <!-- Icon Container -->
        <div class="relative mb-6">
          <div class="aspect-[9/16] rounded-xl overflow-hidden relative">
            <!-- Gradient background -->
            <div class="absolute inset-0 bg-gradient-to-br from-purple-950/80 via-slate-900/90 to-slate-950"></div>

            <!-- Radial glow -->
            <div class="absolute inset-0 opacity-30" style="background-image: radial-gradient(circle at 50% 50%, rgba(168, 85, 247, 0.4) 0%, transparent 50%);"></div>

            <!-- Circular Avatar with Verified Check and Floating Skill Tag -->
            <div class="absolute inset-0 flex items-center justify-center p-3">
              <svg viewBox="0 0 80 80" class="w-full h-full drop-shadow-2xl" fill="none" xmlns="http://www.w3.org/2000/svg">
                <defs>
                  <linearGradient id="avatarGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#a855f7" />
                    <stop offset="100%" stop-color="#6366f1" />
                  </linearGradient>
                  <linearGradient id="skillTagGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" stop-color="#22c55e" />
                    <stop offset="100%" stop-color="#4ade80" />
                  </linearGradient>
                </defs>

                <!-- Background Circle -->
                <circle cx="40" cy="40" r="28" fill="rgba(168,85,247,0.1)" stroke="rgba(168,85,247,0.3)" stroke-width="1" />

                <!-- Avatar Frame -->
                <circle cx="40" cy="38" r="18" fill="url(#avatarGrad)" opacity="0.9" />

                <!-- Avatar Silhouette (simplified person) -->
                <circle cx="40" cy="32" r="6" fill="rgba(255,255,255,0.3)" />
                <path d="M32 48 Q40 52 48 48 L48 52 Q40 58 32 52 Z" fill="rgba(255,255,255,0.25)" />

                <!-- Verified Check Badge -->
                <g>
                  <circle cx="52" cy="28" r="7" fill="#22c55e" class="animate-pulse" style="animation-duration: 2s;" />
                  <path d="M48 28 L51 31 L56 26" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                </g>

                <!-- Floating Skill Tag -->
                <g class="animate-bounce" style="animation-duration: 3s;">
                  <rect x="18" y="54" width="28" height="10" rx="5" fill="url(#skillTagGrad)" />
                  <text x="32" y="61" font-size="5" fill="white" font-weight="600" text-anchor="middle">鋼琴</text>
                  <!-- Skill icon -->
                  <circle cx="23" cy="59" r="2" fill="white" opacity="0.9" />
                </g>

                <!-- Match Score Badge -->
                <circle cx="58" cy="48" r="8" fill="rgba(99,102,241,0.2)" stroke="rgba(99,102,241,0.5)" stroke-width="1" />
                <text x="58" y="51" font-size="6" fill="#a855f7" font-weight="700" text-anchor="middle">96</text>
                <text x="58" y="55" font-size="3" fill="#a855f7" text-anchor="middle">%</text>

                <!-- Connection Lines (matching) -->
                <g class="opacity-20">
                  <line x1="40" y1="20" x2="40" y2="12" stroke="#a855f7" stroke-width="1" stroke-dasharray="2 1" />
                  <circle cx="40" cy="10" r="2" fill="#a855f7" />
                </g>
              </svg>
            </div>
          </div>
        </div>

        <!-- Content -->
        <h3 class="text-xl font-semibold text-slate-100 text-center mb-2">
          {{ $t('landing.features.matching.title') }}
        </h3>
        <p class="text-sm text-slate-400 text-center leading-relaxed">
          {{ $t('landing.features.matching.description') }}
        </p>

        <!-- Hover glow effect -->
        <div class="absolute inset-0 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"
             style="box-shadow: inset 0 0 40px rgba(168, 85, 247, 0.1);">
        </div>
      </div>

      <!-- Feature 3: Real-time Notifications -->
      <div
        class="group glass-card p-6 lg:p-8 rounded-2xl transition-all duration-500 hover:shadow-rose-500/10 hover:-translate-y-2 opacity-0 animate-fade-in-up"
        style="animation-delay: 300ms;"
      >
        <!-- Icon Container -->
        <div class="relative mb-6">
          <div class="aspect-[9/16] rounded-xl overflow-hidden relative">
            <!-- Gradient background -->
            <div class="absolute inset-0 bg-gradient-to-br from-rose-950/80 via-slate-900/90 to-slate-950"></div>

            <!-- Radial glow -->
            <div class="absolute inset-0 opacity-30" style="background-image: radial-gradient(circle at 50% 50%, rgba(244, 63, 94, 0.4) 0%, transparent 50%);"></div>

            <!-- Smartphone Silhouette with LINE Message Bubble -->
            <div class="absolute inset-0 flex items-center justify-center p-3">
              <svg viewBox="0 0 80 80" class="w-full h-full drop-shadow-2xl" fill="none" xmlns="http://www.w3.org/2000/svg">
                <defs>
                  <linearGradient id="phoneGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#f43f5e" />
                    <stop offset="100%" stop-color="#fb7185" />
                  </linearGradient>
                  <linearGradient id="lineBubbleGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#06c755" />
                    <stop offset="100%" stop-color="#00b340" />
                  </linearGradient>
                </defs>

                <!-- Smartphone Frame -->
                <rect x="22" y="10" width="36" height="60" rx="5" fill="rgba(244,63,94,0.1)" stroke="rgba(244,63,94,0.4)" stroke-width="1.5" />

                <!-- Phone Screen -->
                <rect x="25" y="15" width="30" height="48" rx="2" fill="rgba(0,0,0,0.3)" />

                <!-- Status Bar -->
                <rect x="25" y="15" width="30" height="6" rx="1" fill="rgba(244,63,94,0.2)" />
                <circle cx="48" cy="18" r="2" fill="rgba(255,255,255,0.5)" />

                <!-- Chat Header -->
                <rect x="25" y="21" width="30" height="8" rx="1" fill="rgba(244,63,94,0.3)" />
                <circle cx="30" cy="25" r="2.5" fill="rgba(255,255,255,0.3)" />
                <rect x="34" y="23" width="15" height="3" rx="1" fill="rgba(255,255,255,0.4)" />

                <!-- LINE Message Bubble -->
                <g>
                  <!-- Bubble tail -->
                  <path d="M55 38 L62 42 L55 44 Z" fill="url(#lineBubbleGrad)" />

                  <!-- Main bubble -->
                  <rect x="38" y="36" width="24" height="12" rx="4" fill="url(#lineBubbleGrad)" />

                  <!-- LINE logo mark inside bubble -->
                  <circle cx="44" cy="42" r="3" fill="white" opacity="0.9" />
                  <path d="M43 42 L44.5 44 L47 41" stroke="#06c755" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                </g>

                <!-- Notification Badge -->
                <g>
                  <circle cx="54" cy="18" r="4" fill="#ef4444" class="animate-pulse">
                    <animate attributeName="r" values="4;5;4" dur="1s" repeatCount="indefinite" />
                  </circle>
                  <text x="54" y="20" font-size="4" fill="white" font-weight="700" text-anchor="middle">1</text>
                </g>

                <!-- Second notification preview -->
                <rect x="27" y="52" width="20" height="6" rx="2" fill="rgba(255,255,255,0.1)" />
                <circle cx="30" cy="55" r="1.5" fill="rgba(255,255,255,0.3)" />
                <rect x="33" y="54" width="12" height="2" rx="1" fill="rgba(255,255,255,0.2)" />

                <!-- Home indicator -->
                <rect x="32" y="60" width="16" height="2" rx="1" fill="rgba(244,63,94,0.5)" />

                <!-- Notification ring effect -->
                <circle cx="55" cy="16" r="8" fill="none" stroke="#f43f5e" stroke-width="0.5" opacity="0.5" class="animate-ping" style="animation-duration: 1.5s;" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Content -->
        <h3 class="text-xl font-semibold text-slate-100 text-center mb-2">
          {{ $t('landing.features.notifications.title') }}
        </h3>
        <p class="text-sm text-slate-400 text-center leading-relaxed">
          {{ $t('landing.features.notifications.description') }}
        </p>

        <!-- Hover glow effect -->
        <div class="absolute inset-0 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"
             style="box-shadow: inset 0 0 40px rgba(244, 63, 94, 0.1);">
        </div>
      </div>
    </div>

    <!-- CTA -->
    <div class="text-center mt-16">
      <button
        @click="$emit('line-login')"
        :disabled="loading"
        class="group relative inline-flex items-center justify-center gap-3 px-10 py-4 text-lg font-medium text-white transition-all duration-300 rounded-xl overflow-hidden"
      >
        <!-- Gradient Background -->
        <div class="absolute inset-0 bg-gradient-to-r from-indigo-600 to-purple-600 transition-all duration-300 group-hover:from-indigo-500 group-hover:to-purple-500"></div>

        <!-- Glow Effect -->
        <div class="absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300" style="box-shadow: 0 0 60px rgba(99, 102, 241, 0.4);"></div>

        <!-- Content -->
        <span class="relative z-10 flex items-center gap-3">
          <svg v-if="loading" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v16a8 8 0 00-8 8h16a8 8 0 018-8z" />
          </svg>
          <svg v-else class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.477 2 2 6.477 2 12c0 1.891.5 3.657 1.363 5.203L2 22l4.797-1.363A9.956 9.956 0 0 0 12 22c5.523 0 10-4.477 10-10S17.523 2 12 2z"/>
          </svg>
          <span>{{ loading ? $t('landing.cta.connecting') : $t('landing.cta.start_using') }}</span>
        </span>

        <!-- Hover Arrow -->
        <svg class="w-5 h-5 relative z-10 transition-transform duration-300 group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
        </svg>
      </button>
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
/* Glass Card - Unified premium styling */
.glass-card {
  @apply bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl shadow-black/20;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.glass-card:hover {
  @apply border-white/20;
  box-shadow: 0 25px 50px -12px rgba(99, 102, 241, 0.15);
}

/* Staggered Entry Animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in-up {
  animation: fadeInUp 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

/* Pulse animation for notification dot */
@keyframes pulse-glow {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.7;
    transform: scale(1.1);
  }
}

.animate-pulse {
  animation: pulse-glow 2s ease-in-out infinite;
}

/* CTA Button styles */
button {
  @apply hover:scale-[1.02] hover:shadow-indigo-500/20;
}

button:active {
  @apply scale-[0.98];
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .glass-card {
    @apply p-5 rounded-xl;
  }
}
</style>
