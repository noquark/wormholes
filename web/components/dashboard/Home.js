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
  return (
    <section className='home'>
      <header>
        <h1>Welcome to Wormholes</h1>
        <ul className='cards'>
          <StatsCard title='504,558' desc='Total Links Created' />
          <StatsCard title='745' desc='Unique Tags' />
          <StatsCard title='10,504,558' desc='Total Clicks' />
          <StatsCard title='464,484' desc='Unique Users' />
          <StatsCard title='3.4 GB' desc='Database Size' />
        </ul>
        <Divider />
      </header>
    </section>
  )
}
