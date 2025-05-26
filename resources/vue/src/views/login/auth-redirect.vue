<script>
export default {
  name: 'AuthRedirect',
  created() {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const state = urlParams.get('state')
    const provider = urlParams.get('provider') || this.extractProviderFromState(state)
    
    if (code) {
      // Store the OAuth data in localStorage
      const oauthData = {
        code: code,
        state: state,
        provider: provider,
        timestamp: Date.now()
      }
      
      window.localStorage.setItem('x-admin-oauth-data', JSON.stringify(oauthData))
      window.close()
    }
  },
  methods: {
    extractProviderFromState(state) {
      if (!state) return ''
      // State format: provider_timestamp
      const parts = state.split('_')
      if (parts.length >= 1) {
        return parts[0]
      }
      return ''
    }
  },
  render: function (h) {
    return h() // avoid warning message
  }
}
</script>
