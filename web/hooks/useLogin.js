import { useReducer } from 'react'
import { isEmail, normalizeEmail } from 'validator'

const ERR_INVALID_EMAIL = 'Please enter valid email'
const ERR_EMPTY_EMAIL = "Password can't be empty"
const ERR_SHORT_PASSWORD = 'Password is too short'
const ERR_EMPTY_PASSWORD = "Password can't be empty"
const ERR_INVALID_ACTION = "Don't know what to do with given data"
const IS_VALID = ''
const MIN_PASSWORD_LENGTH = 8

const initialState = {
  loggedIn: false,
  email: '',
  password: '',
  emailError: '',
  passwordError: '',
}

function loginReducer(state, action) {
  if (
    typeof action === 'object' ||
    Object.keys(action).every((key) => Object.keys(initialState).includes(key))
  ) {
    return {
      ...state,
      ...action,
    }
  } else {
    console.error(ERR_INVALID_ACTION)
  }
  return state
}

export default function useLogin() {
  const [state, dispatch] = useReducer(loginReducer, initialState)

  function handleEmail(e) {
    const email = e?.target?.value
    if (isEmail(email)) {
      dispatch({ email: normalizeEmail(email), emailError: IS_VALID })
    } else {
      dispatch({ email, emailError: ERR_INVALID_EMAIL })
    }
  }

  function handlePassword(e) {
    const password = e?.target?.value
    if (password?.length && password?.length > MIN_PASSWORD_LENGTH) {
      dispatch({ password, passwordError: IS_VALID })
    } else {
      dispatch({ password, passwordError: ERR_SHORT_PASSWORD })
    }
  }

  async function handleSubmit(e) {
    e.preventDefault()
    if (!state.email) {
      dispatch({ emailError: ERR_EMPTY_EMAIL })
      return
    }
    if (!state.password) {
      dispatch({ passwordError: ERR_EMPTY_PASSWORD })
      return
    }
    if (state.emailError || state.passwordError) {
      return
    }
    try {
      const res = await window.fetch('/api/auth/login', {
        headers: {
          Authorization: `Basic ${window.btoa(
            `${state.email}:${state.password}`
          )}`,
        },
      })
      if (res.ok) {
        dispatch({ loggedIn: true })
      }
    } catch (e) {
      console.error(e)
    }
  }

  function classFor(type) {
    switch (type) {
      case 'email':
        return state.email && (state.emailError ? 'is-invalid' : 'is-valid')
      case 'password':
        return (
          state.password && (state.passwordError ? 'is-invalid' : 'is-valid')
        )
      default:
        return ''
    }
  }

  return {
    state,
    handleEmail,
    handlePassword,
    handleSubmit,
    classFor,
  }
}
