<script lang="ts" setup>
import EventTable from "components/EventTable.vue";

const handleRowClick = (rowId) => {
  console.log("Row clicked in parent component with ID:", rowId);
  // Add your handling logic here
};
</script>

<template>
  <q-page class="flex column no-wrap">
    <div
      class="col q-px-md q-pb-md"
      style="display: flex; flex-direction: column; overflow: hidden"
    >
      <EventTable class="full-height" @row-click="handleRowClick" />
    </div>
  </q-page>
</template>

<style>
.q-page {
  height: 100vh;
  min-height: 0;
}

.full-height {
  height: 100%;
}
</style>
