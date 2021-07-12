import Input from '@/components/Input'
import Layout from '@/components/Layout'
import useLogin from '@/hooks/useLogin'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

export default function LoginPage() {
  const { state, handleEmail, handlePassword, handleSubmit } = useLogin()
  const router = useRouter()

  useEffect(() => {
    if (state.loggedIn) router.push('/')
  }, [state.loggedIn])

  return (
    <Layout>
      <form className='sign-in' onSubmit={handleSubmit}>
        <h1>Wormholes</h1>
        <Input
          type='email'
          label='Email address'
          value={state.email}
          error={state.emailError}
          onChange={handleEmail}
          placeholder='name@example.com'
        />
        <Input
          type='password'
          label='Password'
          value={state.password}
          error={state.passwordError}
          onChange={handlePassword}
          placeholder='password'
        />
        <button type='submit'>Get Started</button>
      </form>
    </Layout>
  )
}
