export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  
  // In development, use localhost:8888 directly
  // In production, use the environment variable
  const apiBase = process.env.API_BASE_URL || 'http://localhost:8888/api/v1'
  const path = getRouterParam(event, 'path')

  const url = `${apiBase}/${path}`

  try {
    const response = await fetch(url, {
      method: event.method,
      headers: {
        ...getRequestHeaders(event),
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: event.method !== 'GET' && event.method !== 'HEAD' ? JSON.stringify(await readBody(event)) : undefined,
    })

    const data = await response.json()

    return data
  } catch (error) {
    return {
      code: 500,
      message: 'API request failed',
      error: error instanceof Error ? error.message : 'Unknown error',
    }
  }
})
