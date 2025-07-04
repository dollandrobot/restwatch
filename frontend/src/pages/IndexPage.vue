<script lang="ts" setup>
import { ref, onMounted } from "vue";
import EventTable from "components/EventTable.vue";
import { GetUserOptions } from "../../wailsjs/go/main/App";
import type { main } from "../../wailsjs/go/models";

const userOptions = ref<main.UserOptions>();
const rightDrawerOpen = ref(false);
const selectedMessage = ref<main.Message>();

const handleRowClick = (row: main.Message) => {
  console.log(row);
  selectedMessage.value = row;
  rightDrawerOpen.value = true;
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

  <q-drawer
    v-model="rightDrawerOpen"
    side="right"
    overlay
    bordered
    :width="800"
  >
    <div class="row q-px-md q-py-sm q-pb-md">
      <div class="col-grow text-h6">Event Details</div>
      <div>
        <q-btn icon="close" flat round dense @click="rightDrawerOpen = false" />
      </div>
      <div></div>
    </div>
    <div class="q-px-md q-py-sm q-pb-md">
      <div class="text-xs text-uppercase text-bold">Body</div>

      <q-markdown
        :content-style="{
          backgroundColor: '#f5f5f5',
          borderRadius: '5px',
        }"
        :src="selectedMessage?.bodyMarkdown"
      />

      <div class="text-xs text-uppercase text-bold">Method</div>
      <div class="q-pa-sm q-mb-sm data rounded-borders">
        {{ selectedMessage?.method }}
      </div>
      <div class="text-xs text-uppercase text-bold">Content Length</div>
      <div class="q-pa-sm q-mb-sm data rounded-borders">
        {{ selectedMessage?.contentLength }}
      </div>
      <div class="text-xs text-uppercase text-bold">Remote Address</div>
      <div class="q-pa-sm q-mb-sm data rounded-borders">
        {{ selectedMessage?.remoteAddr }}
      </div>
      <div class="text-xs text-uppercase text-bold">Received At</div>
      <div class="q-pa-sm q-mb-sm data rounded-borders">
        {{ new Date(selectedMessage?.receivedAt).toLocaleString() }}
      </div>
    </div>
  </q-drawer>
</template>

<style lang="sass" scoped>
.q-page
  height: 100vh
  min-height: 0

.full-height
  height: 100%

.data
  background-color: rgb(0, 0, 0, 0.1)

.body--dark .data
  background-color: rgba(255, 255, 255, 0.1)
</style>
