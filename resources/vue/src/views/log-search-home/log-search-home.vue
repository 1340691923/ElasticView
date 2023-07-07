<template>
  <div class="app-container">
    <log-search value="gn"></log-search>
  </div>
</template>

<script>
export default {
  name: "LogSearchHome",
  components: {
    LogSearch: () => import("@/components/LogSearch/index")
  },
  data() {
    return {};
  },
  computed: {},
  mounted() {},
  methods: {}
};
</script>

<style></style>
