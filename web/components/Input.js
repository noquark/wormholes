import clsx from 'clsx'

export default function Input({ type, label, value, error, ...props }) {
  return (
    <div className='form'>
      <input
        type={type}
        value={value}
        className={clsx([value && (error ? 'is-invalid' : 'is-valid')])}
        id={type}
        {...props}
      />
      <label htmlFor={type}>{label}</label>
      {error && <div className='invalid-tooltip'>{error}</div>}
    </div>
  )
}
