<script lang="ts" setup>
import { ref, onMounted } from "vue";
import EventTable from "components/EventTable.vue";
import { GetUserOptions } from "../../wailsjs/go/main/App";
import type { main } from "../../wailsjs/go/models";

const userOptions = ref<main.UserOptions>();
const rightDrawerOpen = ref(false);
const selectedMessage = ref<main.Message>();
const drawerWidth = ref(300);

let originalWidth: number;
let originalLeft: number;

const handleRowClick = (row: main.Message) => {
  selectedMessage.value = row;
  rightDrawerOpen.value = true;
};

const handlePan = ({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  evt,
  ...newInfo
}: {
  evt: Event;
  isFirst?: boolean;
  position: { left: number };
}) => {
  if (newInfo.isFirst) {
    originalWidth = drawerWidth.value;
    originalLeft = newInfo.position.left;
  } else {
    const newDelta = newInfo.position.left - originalLeft;
    // Should add (instead of subtract) for left drawer
    const newWidth = Math.max(300, originalWidth - newDelta);
    drawerWidth.value = newWidth;
  }
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
    bordered
    :width="drawerWidth"
  >
    <div class="row fit">
      <div
        style="width: 6px; cursor: col-resize"
        v-touch-pan.horizontal.prevent.mouse.preserveCursor="handlePan"
      ></div>
      <q-scroll-area class="fit col">
        <div class="row q-px-md q-py-sm q-pb-md">
          <div class="col-grow text-h6">
            Event Details #{{ selectedMessage?.number }}
          </div>
          <div>
            <q-btn
              icon="close"
              flat
              round
              dense
              @click="rightDrawerOpen = false"
            />
          </div>
          <div></div>
        </div>
        <div class="q-px-md q-py-sm q-pb-md">
          <div class="text-xs text-uppercase text-bold">Received At</div>
          <div class="q-pa-sm q-mb-sm data rounded-borders">
            {{ new Date(selectedMessage?.receivedAt).toLocaleString() }}
          </div>
          <div class="text-xs text-uppercase text-bold">Body</div>
          <div class="q-pa-sm q-mb-sm data rounded-borders">
            <pre class="scrollable-pre">{{
              selectedMessage?.formattedBody
            }}</pre>
          </div>
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
        </div>
      </q-scroll-area>
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
  background-color: rgba(0, 0, 0, 0.1)

.body--dark .data
  background-color: rgba(255, 255, 255, 0.1)

.scrollable-pre
  white-space: pre-wrap
  overflow-x: auto
  max-width: 100%
</style>
