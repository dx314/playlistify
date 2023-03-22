import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'
import {AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE} from "./config";

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
      <App />
  </React.StrictMode>,
)
