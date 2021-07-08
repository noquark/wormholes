import Layout from '@/components/Layout'

export default function IndexPage() {
  return (
    <Layout>
      <form className='sign-in text-center'>
        <h1 className='h1 fw-bold'>Wormholes</h1>
        <div className='form-floating'>
          <input
            type='email'
            className='form-control'
            id='floatingInput'
            placeholder='name@example.com'
          />
          <label htmlFor='floatingInput'>Email address</label>
        </div>
        <div className='form-floating'>
          <input
            type='password'
            className='form-control'
            id='floatingPassword'
            placeholder='Password'
          />
          <label htmlFor='floatingPassword'>Password</label>
        </div>

        <button className='mt-2 w-100 btn btn-lg btn-primary' type='submit'>
          Get Started
        </button>
      </form>
    </Layout>
  )
}
