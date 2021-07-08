import Layout from '@/components/Layout'
import useLogin from '@/hooks/useLogin'
import clsx from 'clsx'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

export default function LoginPage() {
  const { state, handleEmail, handlePassword, handleSubmit, classFor } =
    useLogin()
  const router = useRouter()

  useEffect(() => {
    if (state.loggedIn) router.push('/')
  }, [state.loggedIn])

  return (
    <Layout>
      <form className='sign-in' onSubmit={handleSubmit}>
        <h1 className='h1 fw-bold'>Wormholes</h1>
        <div className='form-floating'>
          <input
            type='email'
            value={state.email}
            className={clsx(['form-control', classFor('email')])}
            id='floatingInput'
            onChange={handleEmail}
            placeholder='name@example.com'
          />
          <label htmlFor='floatingInput'>Email address</label>
          {state.emailError && (
            <div className='invalid-tooltip'>{state.emailError}</div>
          )}
        </div>

        <div className='form-floating'>
          <input
            type='password'
            value={state.password}
            className={clsx(['form-control', classFor('password')])}
            id='floatingPassword'
            onChange={handlePassword}
            placeholder='Password'
          />
          <label htmlFor='floatingPassword'>Password</label>
          {state.passwordError && (
            <div className='invalid-tooltip'>{state.passwordError}</div>
          )}
        </div>

        <button className='mt-3 w-100 btn btn-lg btn-primary' type='submit'>
          Get Started
        </button>
      </form>
    </Layout>
  )
}
