import Head from 'next/head'

export default function Layout({
  children,
  title = 'Wormholes : A self-hosted link shortener',
  description = 'Wormholes is a fast and fail-safe link shortener',
  url = '',
}) {
  return (
    <>
      <Head>
        <title>{title}</title>
        <meta name='description' content={description} />
      </Head>
      <main className='flex flex-shrink-0 bg-light'>{children}</main>
      <footer className='footer mt-auto py-3 text-center bg-light'>
        <p className='mt-5 mb-3 text-muted small'>&copy; 2021 Mohit Singh.</p>
      </footer>
    </>
  )
}
