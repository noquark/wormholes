import Head from 'next/head'

export default function Layout({
  children,
  title = 'Wormholes : A self-hosted link shortener',
  description = 'Wormholes is a fast and fail-safe link shortener',
  full = false,
}) {
  return (
    <>
      <Head>
        <title>{title}</title>
        <meta name='description' content={description} />
      </Head>
      <main className={full && 'full'}>{children}</main>
      {!full && (
        <footer>
          <p>&copy; 2021 Mohit Singh.</p>
        </footer>
      )}
    </>
  )
}
