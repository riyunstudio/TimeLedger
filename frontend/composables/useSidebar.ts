const sidebarState = {
  isOpen: ref(false),
}

export const useSidebar = () => {
  return {
    isOpen: sidebarState.isOpen,
    toggle: () => {
      sidebarState.isOpen.value = !sidebarState.isOpen.value
    },
    open: () => {
      sidebarState.isOpen.value = true
    },
    close: () => {
      sidebarState.isOpen.value = false
    },
  }
}
