import Layout from '@/components/Layout'
import Sidebar from '@/components/Sidebar'
import useUser, { Status } from '@/hooks/useUser'
import { useRouter } from 'next/router'

export default function IndexPage() {
  const [user, status] = useUser()
  const router = useRouter()

  if (status === Status.SUCCESS) {
    return (
      <Layout full>
        <Sidebar />
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
