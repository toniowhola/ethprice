<template>
<div id="page-prices">
    <div class='container'>
      <h1>Prices - ETH in USD</h1>
      <p v-if="$fetchState.pending">Fetching posts...</p>
      <p v-else-if="$fetchState.error">
        Error while fetching posts: {{ $fetchState.error.message }}
      </p>
      <div v-else>
        <el-table
          :data="prices"
          style="width: 100%">
          <el-table-column
            prop="id"
            label="ID"
            width="180">
          </el-table-column>
          <el-table-column
            prop="coin"
            label="Coin"
            width="180">
          </el-table-column>
          <el-table-column
            prop="price"
            label="Price">
          </el-table-column>
          <el-table-column
            prop="created_at"
            label="created_at">
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EthPrices',
  data() {
    return {
      prices: []
    }
  },
  async fetch() {
    this.prices = await this.$axios.$get(
      'http://localhost:10000/prices'
    )
  }
}
</script>

<style>
  .container {
    max-width: 1000px;
    margin: 60px auto;
  }
</style>
