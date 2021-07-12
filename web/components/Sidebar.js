import { useRouter } from 'next/router'
import { Home, LogOut, PieChart, Settings, Users } from 'react-feather'

export function SideLink({ active, icon: Icon, ...props }) {
  return (
    <li>
      <a className={active && 'active'} aria-current='page' {...props}>
        <Icon size={24} />
      </a>
    </li>
  )
}

export default function Sidebar() {
  const router = useRouter()

  async function handleLogout() {
    try {
      const res = await window.fetch('/api/auth/logout')
      if (res.ok) {
        router.push('/login')
      }
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <div className='sidebar'>
      <ul>
        <SideLink title='Home' icon={Home} />
        <SideLink href='#' title='Dashboard' icon={PieChart} />
        <SideLink href='#' title='Users' icon={Users} />
      </ul>
      <ul>
        <SideLink href='#' title='Settings' icon={Settings} />
        <SideLink
          href='#'
          title='Logout'
          onClick={handleLogout}
          icon={LogOut}
        />
      </ul>
    </div>
  )
}
