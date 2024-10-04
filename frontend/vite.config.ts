import basicSsl from '@vitejs/plugin-basic-ssl'

export default {
  base: process.env.BASE_URL,
  plugins: [
    basicSsl()
  ]
}