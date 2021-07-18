import Home from '@/components/dashboard/Home'
import Layout from '@/components/Layout'
import Sidebar from '@/components/Sidebar'
import useData, { Status } from '@/hooks/useData'
import { useRouter } from 'next/router'

export default function IndexPage() {
  const [user, status] = useData('api/auth/user')
  const router = useRouter()

  if (status === Status.SUCCESS) {
    return (
      <Layout full>
        <Sidebar />
        <Home />
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
