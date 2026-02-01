// 共享的資源緩存，避免重複 API 調用
const resourceCache = ref<{
  offerings: any[]
  teachers: Map<number, any>
  rooms: Map<number, any>
  loaded: boolean
}>({
  offerings: [],
  teachers: new Map(),
  rooms: new Map(),
  loaded: false,
})

// 使用 ref 來確保在組件間共享同一個 Promise
const fetchPromise = ref<Promise<void> | null>(null)

export function useResourceCache() {
  const { getCenterId } = useCenterId()
  const authStore = useAuthStore()

  const isAdmin = computed(() => authStore.isAdmin)

  const fetchAllResources = async () => {
    // 如果已經在加載中，返回現有的 promise
    if (fetchPromise.value) {
      return fetchPromise.value
    }

    // 如果已經加載過，直接返回
    if (resourceCache.value.loaded) {
      return
    }

    // 老師端不需要調用 Admin API，直接跳過
    // 老師端的課表數據已經包含足夠的資訊（room_id, title 等）
    if (!isAdmin.value) {
      resourceCache.value.loaded = true
      return
    }

    // 創建新的 promise
    const promise = (async () => {
      try {
        const api = useApi()
        const centerId = getCenterId()

        const [offeringsRes, teachersRes, roomsRes] = await Promise.all([
          // offerings 使用 active endpoint 獲取所有可用的班別
          api.get<any>(`/admin/offerings/active`),
          // teachers 使用 /admin/teachers 端點（useApi 會自動添加 /api/v1 前綴）
          api.get<any[]>(`/admin/teachers`),
          // rooms 使用 /admin/rooms/active 端點
          api.get<any[]>(`/admin/rooms/active`)
        ])

        // 處理 offerings - GetActiveOfferings 返回 []models.Offering 陣列 (使用 data 欄位)
        resourceCache.value.offerings = offeringsRes?.data || offeringsRes || []

        // 處理 teachers - ListTeachers 返回 []AdminTeacherResponse 陣列 (使用 data 欄位)
        const teachersData = (teachersRes as any)?.data || (teachersRes as any)?.datas || teachersRes || []
        resourceCache.value.teachers = new Map()
        if (Array.isArray(teachersData)) {
          teachersData.forEach((t: any) => {
            resourceCache.value.teachers.set(t.id, t)
          })
        }

        // 處理 rooms - GetRooms 返回 []RoomResponse 陣列 (使用 data 欄位)
        const roomsData = (roomsRes as any)?.data || (roomsRes as any)?.datas || roomsRes || []
        resourceCache.value.rooms = new Map()
        if (Array.isArray(roomsData)) {
          roomsData.forEach((r: any) => {
            resourceCache.value.rooms.set(r.id, r)
          })
        }

        resourceCache.value.loaded = true
      } catch (error) {
        console.error('Failed to fetch resources:', error)
        resourceCache.value.loaded = false
        fetchPromise.value = null
        throw error
      }
    })()

    fetchPromise.value = promise
    return promise
  }

  const getTeacherName = (teacherId: number): string => {
    const teacher = resourceCache.value.teachers.get(teacherId)
    return teacher?.name || `老師 ${teacherId}`
  }

  const getRoomName = (roomId: number): string => {
    const room = resourceCache.value.rooms.get(roomId)
    return room?.name || `教室 ${roomId}`
  }

  const getOfferingName = (offeringId: number): string => {
    const offering = resourceCache.value.offerings.find(o => o.id === offeringId)
    return offering?.name || `課程 ${offeringId}`
  }

  const invalidate = (type?: string) => {
    // 雖然目前 fetchAllResources 是全量抓取，但我們保留 type 參數以利未來擴充
    resourceCache.value.loaded = false
    fetchPromise.value = null
  }

  const fetchIfExpired = async () => {
    return fetchAllResources()
  }

  return {
    resourceCache,
    fetchAllResources,
    getTeacherName,
    getRoomName,
    getOfferingName,
    invalidate,
    fetchIfExpired,
  }
}
