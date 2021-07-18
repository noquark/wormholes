import { useEffect, useState } from 'react'

export const Status = {
  LOADING: 'loading',
  ERROR: 'error',
  SUCCESS: 'success',
}

export default function useData(endpoint) {
  const [data, setData] = useState()
  const [status, setStatus] = useState(Status.LOADING)

  useEffect(() => {
    async function fetchData() {
      try {
        const res = await window.fetch(endpoint)
        if (res.ok) {
          const data = await res.json()
          setData(data)
          setStatus(Status.SUCCESS)
        } else {
          setStatus(Status.ERROR)
        }
      } catch (e) {
        setStatus(Status.ERROR)
        console.error(e)
      }
    }
    fetchData()
  }, [endpoint])
  return [data, status]
}
