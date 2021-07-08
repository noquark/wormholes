import Layout from '@/components/Layout'
import useUser, { Status } from '@/hooks/useUser'
import { useRouter } from 'next/router'

export default function IndexPage() {
  const [user, status] = useUser()
  const router = useRouter()

  if (status === Status.SUCCESS) {
    return (
      <Layout>
        <h1>Welcome to Wormholes !</h1>
        <p>
          Your email is <strong>{user.email}</strong>
        </p>
      </Layout>
    )
  }
  if (status === Status.ERROR) {
    router.push('/login')
  }
  return (
    <Layout>
      <h1>Loading...</h1>
    </Layout>
  )
}
