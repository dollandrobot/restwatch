<script lang="ts" setup>
import { ref, onMounted } from "vue";
import EventTable from "components/EventTable.vue";
import { GetUserOptions } from "../../wailsjs/go/main/App";
import type { main } from "../../wailsjs/go/models";

const userOptions = ref<main.UserOptions>();

const handleRowClick = (rowId: string) => {
  console.log("Row clicked in parent component with ID:", rowId);
  // Add your handling logic here
};

onMounted(async () => {
  try {
    userOptions.value = await GetUserOptions();
  } catch (error) {
    console.log("error getting user options", error);
  }
});
</script>

<template>
  <q-page class="flex column no-wrap">
    <div
      class="col q-px-md q-pb-md"
      style="display: flex; flex-direction: column; overflow: hidden"
    >
      <EventTable v-if="userOptions" :userOptions @row-click="handleRowClick" />
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
