<template>
  <div class="relative min-h-screen overflow-hidden">
    <!-- Unified Animated Mesh Gradient Background -->
    <div class="fixed inset-0 z-0">
      <!-- Base gradient -->
      <div class="absolute inset-0 bg-gradient-to-br from-slate-950 via-indigo-950 to-slate-950"></div>
      
      <!-- Animated Gradient Orbs -->
      <div class="absolute inset-0 animate-gradient-slow">
        <div class="gradient-orb-1"></div>
        <div class="gradient-orb-2"></div>
        <div class="gradient-orb-3"></div>
        <div class="gradient-orb-4"></div>
      </div>
      
      <!-- Subtle Grid Overlay -->
      <div class="absolute inset-0 pointer-events-none opacity-[0.03]" style="background-image: radial-gradient(circle, #6366f1 1px, transparent 1px); background-size: 40px 40px;"></div>
      
      <!-- Radial Vignette -->
      <div class="absolute inset-0 bg-gradient-to-t from-slate-950/50 via-transparent to-slate-950/30"></div>
    </div>

    <!-- Main Content -->
    <div class="relative z-10">
      <!-- Hero Section -->
      <section class="relative min-h-screen flex items-center">
        <LandingHeroSection :loading="loading" @line-login="handleLineLogin" />
      </section>

      <!-- Demo Section -->
      <section class="relative">
        <LandingDemoSandbox />
      </section>

      <!-- Features Section -->
      <section class="relative pb-20">
        <LandingFeatureShowcase :loading="loading" @line-login="handleLineLogin" />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const router = useRouter()
const loading = ref(false)

const handleLineLogin = async () => {
  loading.value = true

  try {
    router.push('/teacher/login')
  } catch (err) {
    console.error('Login redirect failed:', err)
    loading.value = false
  }
}
</script>

<style>
/* Gradient orb animations - Shared across all landing sections */
@keyframes gradient-shift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.animate-gradient-slow {
  animation: gradient-shift 30s ease-in-out infinite;
}

.gradient-orb-1,
.gradient-orb-2,
.gradient-orb-3,
.gradient-orb-4 {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  animation: float 20s ease-in-out infinite;
}

.gradient-orb-1 {
  top: -30%;
  left: -20%;
  width: 60%;
  height: 60%;
  background: radial-gradient(circle, #6366f1 0%, transparent 70%);
  animation-delay: 0s;
}

.gradient-orb-2 {
  top: 20%;
  right: -30%;
  width: 70%;
  height: 70%;
  background: radial-gradient(circle, #a855f7 0%, transparent 70%);
  animation-delay: -5s;
}

.gradient-orb-3 {
  bottom: -40%;
  left: 20%;
  width: 50%;
  height: 50%;
  background: radial-gradient(circle, #1e3a8a 0%, transparent 70%);
  animation-delay: -10s;
}

.gradient-orb-4 {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
  height: 80%;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.1) 0%, transparent 70%);
  animation: pulse-glow 8s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  25% {
    transform: translate(30px, -30px) scale(1.1);
  }
  50% {
    transform: translate(-20px, 20px) scale(0.95);
  }
  75% {
    transform: translate(20px, 30px) scale(1.05);
  }
}

@keyframes pulse-glow {
  0%, 100% {
    opacity: 0.5;
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    opacity: 0.8;
    transform: translate(-50%, -50%) scale(1.2);
  }
}
</style>
