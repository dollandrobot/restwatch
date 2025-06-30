<script lang="ts">
const columns = [
  {
    name: 'name',
    required: true,
    label: 'Event',
    align: 'left',
    field: row => row.RawMessage,
    format: val => `${val}`,
    sortable: true
  },
  {
    name: 'receivedAt',
    required: true,
    label: 'Received At',
    align: 'left',
    field: row => row.receivedAt,
    format: val => `${val}`,
    sortable: true
  }
]

const rows = []

export default {
  setup () {
    return {
      columns,
      rows
    }
  }
}
</script>


<template>
  <div class="q-pa-md">
    <q-table
        title="Treats"
        :rows="rows"
        :columns="columns"
        row-key="name"
    />
  </div>
</template>

<style scoped>

</style>