import { useEffect, useState } from 'react'

export const Status = {
  LOADING: 'loading',
  ERROR: 'error',
  SUCCESS: 'success',
}

export default function useUser() {
  const [user, setUser] = useState()
  const [status, setStatus] = useState(Status.LOADING)

  useEffect(() => {
    async function fetchUser() {
      try {
        const res = await window.fetch('/api/auth/user')
        if (res.ok) {
          const user = await res.json()
          setUser(user)
          setStatus(Status.SUCCESS)
        } else {
          setStatus(Status.ERROR)
        }
      } catch (e) {
        setStatus(Status.ERROR)
        console.error(e)
      }
    }
    fetchUser()
  }, [])
  return [user, status]
}
