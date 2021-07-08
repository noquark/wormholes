module.exports = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  async rewrites() {
    return [
      {
        source: '/api/:slug*',
        destination: 'http://localhost:5000/api/v1/:slug*',
      },
    ]
  },
}
