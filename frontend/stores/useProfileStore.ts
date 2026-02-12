import { defineStore } from 'pinia'
import type { TeacherSkill, TeacherCertificate, Teacher, Hashtag } from '~/types'
import { withLoading } from '~/utils/loadingHelper'

export interface BackgroundImage {
  id: number
  teacher_id: number
  path?: string  // 舊版相容性欄位
  filename?: string  // 舊版相容性欄位
  file_url: string
  url: string
  created_at: string
}

export const useProfileStore = defineStore('profile', () => {
  // 資料狀態
  const skills = ref<TeacherSkill[]>([])
  const certificates = ref<TeacherCertificate[]>([])
  const profile = ref<Teacher | null>(null)
  const backgrounds = ref<BackgroundImage[]>([])

  // Loading 狀態
  const isFetching = ref(false)
  const isUpdating = ref(false)
  const isCreatingSkill = ref(false)
  const isUpdatingSkill = ref(false)
  const isDeletingSkill = ref(false)
  const isCreatingCertificate = ref(false)
  const isUpdatingCertificate = ref(false)
  const isDeletingCertificate = ref(false)
  const isFetchingBackgrounds = ref(false)
  const isUploadingBackground = ref(false)
  const isDeletingBackground = ref(false)

  // 技能操作
  const fetchSkills = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<TeacherSkill[]>('/teacher/me/skills')
        skills.value = response || []
      } catch (error) {
        console.error('Failed to fetch skills:', error)
        throw error
      }
    })
  }

  const createSkill = async (data: { category: string; skill_name: string; hashtag_ids?: number[] }) => {
    return withLoading(isCreatingSkill, async () => {
      const api = useApi()
      const response = await api.post<TeacherSkill>('/teacher/me/skills', data)
      skills.value.push(response)
      return response
    })
  }

  const updateSkill = async (skillId: number, data: { category: string; skill_name: string; hashtags?: string[] }) => {
    return withLoading(isUpdatingSkill, async () => {
      const api = useApi()
      const response = await api.put<TeacherSkill>(`/teacher/me/skills/${skillId}`, data)
      const index = skills.value.findIndex(s => s.id === skillId)
      if (index !== -1) {
        skills.value[index] = response
      }
      return response
    })
  }

  const deleteSkill = async (skillId: number) => {
    return withLoading(isDeletingSkill, async () => {
      const api = useApi()
      await api.delete(`/teacher/me/skills/${skillId}`)
      skills.value = skills.value.filter(s => s.id !== skillId)
    })
  }

  // 證照操作
  const fetchCertificates = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<TeacherCertificate[]>('/teacher/me/certificates')
        certificates.value = response || []
      } catch (error) {
        console.error('Failed to fetch certificates:', error)
        throw error
      }
    })
  }

  const createCertificate = async (data: {
    name: string
    file_url?: string
    issued_at?: string
    visibility?: number
  }) => {
    return withLoading(isCreatingCertificate, async () => {
      const api = useApi()
      const response = await api.post<TeacherCertificate>('/teacher/me/certificates', data)
      certificates.value.push(response)
      return response
    })
  }

  const updateCertificate = async (certId: number, data: {
    name: string
    file_url?: string
    issued_at?: string
    visibility?: number
  }) => {
    return withLoading(isUpdatingCertificate, async () => {
      const api = useApi()
      const response = await api.put<TeacherCertificate>(`/teacher/me/certificates/${certId}`, data)
      const index = certificates.value.findIndex(c => c.id === certId)
      if (index !== -1) {
        certificates.value[index] = response
      }
      return response
    })
  }

  const deleteCertificate = async (certId: number) => {
    return withLoading(isDeletingCertificate, async () => {
      const api = useApi()
      await api.delete(`/teacher/me/certificates/${certId}`)
      certificates.value = certificates.value.filter(c => c.id !== certId)
    })
  }

  // 個人檔案操作
  const fetchProfile = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<Teacher>('/teacher/me/profile')
        profile.value = response || null
        return response
      } catch (error) {
        console.error('Failed to fetch profile:', error)
        throw error
      }
    })
  }

  const updateProfile = async (data: Partial<Teacher>) => {
    return withLoading(isUpdating, async () => {
      const api = useApi()
      const response = await api.put<Teacher>('/teacher/me/profile', data)
      profile.value = response || null
      return response
    })
  }

  // 標籤相關方法
  const searchHashtags = async (query: string): Promise<Hashtag[]> => {
    if (!query || query.length < 1) {
      return []
    }

    try {
      const api = useApi()
      const response = await api.get<Hashtag[]>('/hashtags/search', { q: query })
      return response || []
    } catch (error) {
      console.error('Failed to search hashtags:', error)
      return []
    }
  }

  const createHashtag = async (name: string): Promise<Hashtag | null> => {
    if (!name || name.length < 2) {
      return null
    }

    try {
      const api = useApi()
      const response = await api.post<Hashtag>('/hashtags', { name: '#' + name })
      return response || null
    } catch (error) {
      console.error('Failed to create hashtag:', error)
      return null
    }
  }

  const processHashtag = async (tagName: string): Promise<number | null> => {
    if (!tagName || tagName.length < 2) {
      return null
    }

    try {
      const searchResults = await searchHashtags(tagName)
      const existing = searchResults.find(h => h.name === tagName || h.name === '#' + tagName)

      if (existing) {
        return existing.id
      }

      const newHashtag = await createHashtag(tagName)
      return newHashtag?.id || null
    } catch (error) {
      console.error('Failed to process hashtag:', tagName, error)
      return null
    }
  }

  // 背景圖片操作
  interface TeacherBackgroundResponse {
    id: number
    teacher_id: number
    file_url: string
    created_at: string
  }

  const fetchBackgrounds = async () => {
    return withLoading(isFetchingBackgrounds, async () => {
      try {
        const api = useApi()
        const response = await api.get<TeacherBackgroundResponse[]>('/teacher/me/backgrounds')
        // 轉換為 BackgroundImage 格式
        backgrounds.value = (response || []).map((bg: TeacherBackgroundResponse) => ({
          id: bg.id,
          teacher_id: bg.teacher_id,
          file_url: bg.file_url,
          url: bg.file_url, // API 已回傳完整 URL
          created_at: bg.created_at,
        }))
      } catch (error) {
        console.error('Failed to fetch backgrounds:', error)
        throw error
      }
    })
  }

  const uploadBackground = async (file: File): Promise<BackgroundImage> => {
    return withLoading(isUploadingBackground, async () => {
      const api = useApi()
      const response = await api.upload<TeacherBackgroundResponse>('/teacher/me/backgrounds', file)

      const newImage: BackgroundImage = {
        id: response.id,
        teacher_id: response.teacher_id,
        file_url: response.file_url,
        url: response.file_url,
        created_at: response.created_at,
      }
      backgrounds.value.push(newImage)
      return newImage
    })
  }

  const deleteBackground = async (id: number) => {
    return withLoading(isDeletingBackground, async () => {
      const api = useApi()
      await api.delete(`/teacher/me/backgrounds/${id}`)
      backgrounds.value = backgrounds.value.filter(b => b.id !== id)
    })
  }

  return {
    // 資料狀態
    skills,
    certificates,
    profile,
    backgrounds,

    // Loading 狀態
    isFetching,
    isUpdating,
    isCreatingSkill,
    isUpdatingSkill,
    isDeletingSkill,
    isCreatingCertificate,
    isUpdatingCertificate,
    isDeletingCertificate,
    isFetchingBackgrounds,
    isUploadingBackground,
    isDeletingBackground,

    // 技能方法
    fetchSkills,
    createSkill,
    updateSkill,
    deleteSkill,

    // 證照方法
    fetchCertificates,
    createCertificate,
    updateCertificate,
    deleteCertificate,

    // 個人檔案方法
    fetchProfile,
    updateProfile,

    // 標籤方法
    searchHashtags,
    createHashtag,
    processHashtag,

    // 背景圖片方法
    fetchBackgrounds,
    uploadBackground,
    deleteBackground,
  }
}, {
  persist: {
    key: 'timeledger-profile',
    paths: ['skills', 'certificates', 'profile'],
  },
})
