import useData, { Status } from '@/hooks/useData'
import Divider from './Divider'

export function StatsCard({ title, desc }) {
  return (
    <li>
      <h2>{title}</h2>
      <p>{desc}</p>
    </li>
  )
}

export default function Home() {
  const [data, status] = useData('api/stats/cards')
  const formatter = new Intl.NumberFormat()
  if (status === Status.SUCCESS) {
    return (
      <section className='home'>
        <header>
          <h1>Welcome to Wormholes</h1>
          <ul className='cards'>
            <StatsCard
              title={formatter.format(data?.overview?.links)}
              desc='Links Created'
            />
            <StatsCard
              title={formatter.format(data?.overview?.tags)}
              desc='Unique Tags'
            />
            <StatsCard
              title={formatter.format(data?.overview?.clicks)}
              desc='Clicks'
            />
            <StatsCard
              title={formatter.format(data?.overview?.users)}
              desc='Unique Users'
            />
            <StatsCard
              title={`${data?.db_size?.links} / ${data?.db_size?.clicks}`}
              desc='Links / Clicks Size'
            />
          </ul>
          <Divider />
        </header>
      </section>
    )
  }

  if (status === Status.ERROR) {
    return (
      <section className='home'>
        <header>
          <h1>Failed to load</h1>
        </header>
      </section>
    )
  }
  return (
    <section className='home'>
      <header>
        <h1>Loading...</h1>
      </header>
    </section>
  )
}
